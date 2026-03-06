"""Generate turn sheet background images for the demo adventure game.

Called by the generate-demo-images shell wrapper.
Uses the OpenAI DALL-E 3 API to produce portrait images (1024x1792),
then resizes them to the target turn sheet dimensions (1240x1754).
"""

import argparse
import io
import os
import sys
import time

import requests
from PIL import Image

TARGET_WIDTH = 1240
TARGET_HEIGHT = 1754
MAX_FILE_SIZE = 1_048_576  # 1 MB
JPEG_QUALITY_START = 92

STYLE_PREFIX = (
    "Oil painting style, moody medieval atmosphere, warm candlelight "
    "and cool stone tones. Portrait orientation, tall format. "
    "No text, no UI elements, no people. "
)

IMAGES = [
    {
        "name": "join-game",
        "prompt": (
            "The exterior of an old medieval abbey at dusk. Heavy weathered "
            "wooden doors set in a crumbling stone facade. Warm candlelight "
            "glows from a narrow gothic window. A windswept hill beneath an "
            "overcast sky. This is the entrance to a mysterious place."
        ),
    },
    {
        "name": "inventory-management",
        "prompt": (
            "A stone shelf alcove in a medieval abbey wall. A worn leather "
            "satchel sits open with items laid out on the stone: a rusty iron "
            "key, a stubby tallow candle, a small silver cross on a chain, a "
            "cracked leather journal, and a coil of hemp rope. Warm "
            "candlelight from the side."
        ),
    },
    {
        "name": "location-grand-staircase",
        "prompt": (
            "A wide stone staircase sweeping upward through a medieval abbey "
            "entrance hall. Iron wall sconces with flickering candles. Dust "
            "motes in pale light from a high narrow gothic window. Beneath "
            "the lowest step, a small weathered wooden door with a dark iron "
            "handle is barely visible."
        ),
    },
    {
        "name": "location-narrow-passage",
        "prompt": (
            "A cramped stone tunnel with rough-hewn walls in a medieval "
            "underground passage. The ceiling is so low you must stoop. A "
            "single candle in a wall bracket casts long shadows ahead. The "
            "passage stretches into darkness."
        ),
    },
    {
        "name": "location-wine-cellar",
        "prompt": (
            "A vaulted medieval cellar with stone arches. Dusty wine bottles "
            "fill sagging wooden shelves on both sides. Cobwebs drape every "
            "surface. Dim amber torchlight from the far end. Deep shadows "
            "between the wine racks."
        ),
    },
    {
        "name": "location-crypt",
        "prompt": (
            "A cold stone crypt beneath a medieval abbey. Ancient sarcophagi "
            "rest in alcoves along the walls. Faded inscriptions are carved "
            "into the stone floor. Pale blue-grey light seeps from an "
            "unknown source. Faint mist near the floor."
        ),
    },
    {
        "name": "location-underground-chapel",
        "prompt": (
            "A forgotten chapel carved from bedrock deep underground beneath "
            "a medieval abbey. Crumbling wooden pews face a cracked stone "
            "altar. A thin beam of faint natural light filters down through "
            "a narrow shaft in the high ceiling far above."
        ),
    },
    {
        "name": "location-flooded-corridor",
        "prompt": (
            "A vaulted stone corridor beneath a medieval abbey with dark "
            "still water covering the floor. Distant torchlight reflects in "
            "amber ripples on the water surface. Drips fall from the ceiling. "
            "The passage stretches into shadow."
        ),
    },
    {
        "name": "location-abbots-study",
        "prompt": (
            "A hidden room behind a false wall in a medieval abbey. A heavy "
            "oak desk buried under scattered manuscripts and yellowed papers. "
            "Old leather-bound books line stone shelves. A single candle stub "
            "burns in a brass holder, casting warm golden light."
        ),
    },
    {
        "name": "location-well-chamber",
        "prompt": (
            "A circular stone chamber built around an ancient well beneath a "
            "medieval abbey. A frayed rope hangs over the well lip into "
            "darkness below. The stones are slick with moisture. Faint light "
            "from a passage illuminates the curved walls."
        ),
    },
    {
        "name": "location-herb-garden",
        "prompt": (
            "An overgrown walled garden behind a medieval abbey, seen from "
            "inside the garden. Tangled herbs push through cracked "
            "flagstones. Crumbling stone walls frame a grey overcast sky. "
            "The only outdoor location with natural daylight, muted tones."
        ),
    },
    {
        "name": "location-bell-tower-vault",
        "prompt": (
            "A small secret vault beneath a medieval bell tower. Iron-bound "
            "wooden chests sit under low stone arches. Dust motes swirl in a "
            "thin shaft of light falling from above through a gap in the "
            "stonework."
        ),
    },
]


MAX_RETRIES = 3
RETRY_DELAY = 10


def generate_image(api_key: str, prompt: str) -> bytes:
    """Call the DALL-E 3 API and return the raw image bytes, with retries."""
    full_prompt = STYLE_PREFIX + prompt

    for attempt in range(1, MAX_RETRIES + 1):
        try:
            resp = requests.post(
                "https://api.openai.com/v1/images/generations",
                headers={
                    "Authorization": f"Bearer {api_key}",
                    "Content-Type": "application/json",
                },
                json={
                    "model": "dall-e-3",
                    "prompt": full_prompt,
                    "n": 1,
                    "size": "1024x1792",
                    "quality": "hd",
                    "response_format": "url",
                },
                timeout=120,
            )
            resp.raise_for_status()
            data = resp.json()
            image_url = data["data"][0]["url"]

            img_resp = requests.get(image_url, timeout=120)
            img_resp.raise_for_status()
            return img_resp.content

        except requests.exceptions.HTTPError as e:
            status = e.response.status_code if e.response is not None else 0
            if attempt < MAX_RETRIES and status in (403, 429, 500, 502, 503):
                print(f"          retry {attempt}/{MAX_RETRIES} after {status}, waiting {RETRY_DELAY}s...")
                time.sleep(RETRY_DELAY)
            else:
                raise


def resize_image(raw_bytes: bytes, width: int, height: int) -> bytes:
    """Resize image to exact dimensions and return as JPEG bytes under MAX_FILE_SIZE."""
    img = Image.open(io.BytesIO(raw_bytes))
    img = img.convert("RGB")
    img = img.resize((width, height), Image.LANCZOS)

    quality = JPEG_QUALITY_START
    while quality >= 50:
        buf = io.BytesIO()
        img.save(buf, format="JPEG", quality=quality, optimize=True)
        if buf.tell() <= MAX_FILE_SIZE:
            return buf.getvalue()
        quality -= 5

    buf = io.BytesIO()
    img.save(buf, format="JPEG", quality=50, optimize=True)
    return buf.getvalue()


def main():
    parser = argparse.ArgumentParser(description="Generate demo turn sheet images")
    parser.add_argument("--api-key", required=True, help="OpenAI API key")
    parser.add_argument("--output-dir", required=True, help="Target directory")
    parser.add_argument("--all", action="store_true", help="Regenerate all images")
    parser.add_argument("--only", type=str, help="Generate only this image (name without .png)")
    args = parser.parse_args()

    images_to_generate = IMAGES
    if args.only:
        images_to_generate = [i for i in IMAGES if i["name"] == args.only]
        if not images_to_generate:
            names = ", ".join(i["name"] for i in IMAGES)
            print(f"ERROR: Unknown image '{args.only}'. Available: {names}")
            sys.exit(1)

    generated = 0
    skipped = 0
    failed = 0

    for entry in images_to_generate:
        output_path = os.path.join(args.output_dir, f"{entry['name']}.jpg")

        if not args.all and not args.only and os.path.exists(output_path):
            file_size = os.path.getsize(output_path)
            if file_size > 10000:
                print(f"  SKIP  {entry['name']}.jpg (already exists, {file_size:,} bytes)")
                skipped += 1
                continue

        print(f"  GEN   {entry['name']}.jpg ...")
        try:
            raw = generate_image(args.api_key, entry["prompt"])
            resized = resize_image(raw, TARGET_WIDTH, TARGET_HEIGHT)
            with open(output_path, "wb") as f:
                f.write(resized)
            file_size = len(resized)
            print(f"  OK    {entry['name']}.jpg ({file_size:,} bytes, q={JPEG_QUALITY_START})")
            generated += 1
        except requests.exceptions.HTTPError as e:
            print(f"  FAIL  {entry['name']}.jpg: {e}")
            if hasattr(e, "response") and e.response is not None:
                print(f"        {e.response.text[:200]}")
            failed += 1
        except Exception as e:
            print(f"  FAIL  {entry['name']}.jpg: {e}")
            failed += 1

        if entry != images_to_generate[-1]:
            time.sleep(1)

    print()
    print(f"Done: {generated} generated, {skipped} skipped, {failed} failed")
    if failed > 0:
        sys.exit(1)


if __name__ == "__main__":
    main()
