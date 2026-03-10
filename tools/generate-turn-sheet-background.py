"""Generate a single turn sheet background image via the OpenAI DALL-E 3 API.

Called by the generate-turn-sheet-background shell wrapper.

Turn sheet backgrounds are portrait-orientation images sized to match A4 at
~150 DPI (1240x1754 by default). DALL-E 3 generates at 1024x1792, and we
resize to the exact target dimensions before saving.

Specifications:
    Orientation : portrait
    DALL-E size : 1024x1792
    Target size : 1240x1754 px (A4 @ ~150 DPI)
    DPI         : ~150 (210mm = 8.27in → 1240/8.27 ≈ 150)
    Format      : JPEG (default) or PNG
    Max file    : 1 MB (JPEG only; PNG is not size-capped)
"""

import argparse
import io
import os
import sys
import time

import requests
from PIL import Image

DEFAULT_WIDTH = 1240
DEFAULT_HEIGHT = 1754
DEFAULT_MAX_SIZE = 1_048_576  # 1 MB
DEFAULT_QUALITY = 92

DEFAULT_STYLE_PREFIX = (
    "Oil painting style, atmospheric fantasy illustration, portrait orientation, "
    "tall format. No text, no UI elements, no people in the foreground. "
)

MAX_RETRIES = 3
RETRY_DELAY = 10


def generate_image(api_key: str, prompt: str) -> bytes:
    """Call the DALL-E 3 API and return the raw image bytes."""
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
                    "prompt": prompt,
                    "n": 1,
                    "size": "1024x1792",
                    "quality": "hd",
                    "response_format": "url",
                },
                timeout=120,
            )
            resp.raise_for_status()
            image_url = resp.json()["data"][0]["url"]

            img_resp = requests.get(image_url, timeout=120)
            img_resp.raise_for_status()
            return img_resp.content

        except requests.exceptions.HTTPError as e:
            status = e.response.status_code if e.response is not None else 0
            if attempt < MAX_RETRIES and status in (403, 429, 500, 502, 503):
                print(f"  retry {attempt}/{MAX_RETRIES} after HTTP {status}, waiting {RETRY_DELAY}s...")
                time.sleep(RETRY_DELAY)
            else:
                raise


def resize_and_save_jpeg(raw_bytes: bytes, width: int, height: int,
                         output_path: str, quality: int, max_size: int) -> int:
    """Resize to exact dimensions and save as JPEG, reducing quality to fit max_size."""
    img = Image.open(io.BytesIO(raw_bytes)).convert("RGB")
    img = img.resize((width, height), Image.LANCZOS)

    q = quality
    while q >= 50:
        buf = io.BytesIO()
        img.save(buf, format="JPEG", quality=q, optimize=True)
        if buf.tell() <= max_size:
            with open(output_path, "wb") as f:
                f.write(buf.getvalue())
            return buf.tell()
        q -= 5

    # Last resort: quality 50
    buf = io.BytesIO()
    img.save(buf, format="JPEG", quality=50, optimize=True)
    with open(output_path, "wb") as f:
        f.write(buf.getvalue())
    return buf.tell()


def resize_and_save_png(raw_bytes: bytes, width: int, height: int,
                        output_path: str) -> int:
    """Resize to exact dimensions and save as PNG."""
    img = Image.open(io.BytesIO(raw_bytes)).convert("RGBA")
    img = img.resize((width, height), Image.LANCZOS)
    img.save(output_path, format="PNG", optimize=True)
    return os.path.getsize(output_path)


def main():
    parser = argparse.ArgumentParser(
        description="Generate a turn sheet background image via DALL-E 3",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""\
examples:
  %(prog)s --name desert-oasis \\
      --prompt "A lush oasis village surrounded by golden sand dunes"

  %(prog)s --name desert-ruins --format png \\
      --prompt "Crumbling ancient desert ruins with hieroglyphs" \\
      --output-dir backend/internal/runner/cli/test_data_images

  %(prog)s --name inventory --style-prefix "Watercolour style, " \\
      --prompt "A treasure room filled with exotic goods"
""",
    )
    parser.add_argument("--api-key", required=True, help=argparse.SUPPRESS)
    parser.add_argument("--default-output-dir", default=".", help=argparse.SUPPRESS)
    parser.add_argument("--name", required=True,
                        help="Output filename without extension")
    parser.add_argument("--prompt", required=True,
                        help="Image generation prompt (style prefix is prepended)")
    parser.add_argument("--output-dir",
                        help="Target directory (default: test_data_images)")
    parser.add_argument("--style-prefix", default=DEFAULT_STYLE_PREFIX,
                        help="Text prepended to the prompt")
    parser.add_argument("--width", type=int, default=DEFAULT_WIDTH,
                        help=f"Target width in pixels (default: {DEFAULT_WIDTH})")
    parser.add_argument("--height", type=int, default=DEFAULT_HEIGHT,
                        help=f"Target height in pixels (default: {DEFAULT_HEIGHT})")
    parser.add_argument("--max-size", type=int, default=DEFAULT_MAX_SIZE,
                        help=f"Max JPEG file size in bytes (default: {DEFAULT_MAX_SIZE})")
    parser.add_argument("--quality", type=int, default=DEFAULT_QUALITY,
                        help=f"Initial JPEG quality 1-100 (default: {DEFAULT_QUALITY})")
    parser.add_argument("--format", choices=["jpg", "png"], default="jpg",
                        help="Output format (default: jpg)")
    parser.add_argument("--force", action="store_true",
                        help="Overwrite existing file")
    args = parser.parse_args()

    output_dir = args.output_dir or args.default_output_dir
    os.makedirs(output_dir, exist_ok=True)

    ext = args.format
    output_path = os.path.join(output_dir, f"{args.name}.{ext}")

    if os.path.exists(output_path) and not args.force:
        size = os.path.getsize(output_path)
        print(f"SKIP  {args.name}.{ext} already exists ({size:,} bytes). Use --force to overwrite.")
        sys.exit(0)

    full_prompt = args.style_prefix + args.prompt
    print(f"GEN   {args.name}.{ext} ...")
    print(f"      prompt: {full_prompt[:120]}...")
    print(f"      target: {args.width}x{args.height} {ext.upper()}")

    try:
        raw = generate_image(args.api_key, full_prompt)

        if ext == "png":
            file_size = resize_and_save_png(raw, args.width, args.height, output_path)
        else:
            file_size = resize_and_save_jpeg(
                raw, args.width, args.height, output_path,
                args.quality, args.max_size,
            )

        print(f"OK    {output_path} ({file_size:,} bytes)")

    except requests.exceptions.HTTPError as e:
        print(f"FAIL  {args.name}.{ext}: {e}")
        if hasattr(e, "response") and e.response is not None:
            print(f"      {e.response.text[:300]}")
        sys.exit(1)
    except Exception as e:
        print(f"FAIL  {args.name}.{ext}: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
