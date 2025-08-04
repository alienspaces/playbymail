/**
 * @typedef {Object} Game
 * @property {string} id
 * @property {string} name
 * @property {string} game_type
 * @property {string} created_at
 * @property {string} [updated_at]
 */

/**
 * @typedef {Object} GameItem
 * @property {string} id
 * @property {string} game_id
 * @property {string} name
 * @property {string} description
 * @property {string} created_at
 * @property {string} [updated_at]
 * @property {string} [deleted_at]
 */

/**
 * @typedef {Object} GameCreature
 * @property {string} id
 * @property {string} game_id
 * @property {string} name
 * @property {string} description
 * @property {string} created_at
 * @property {string} [updated_at]
 * @property {string} [deleted_at]
 */

/**
 * @typedef {Object} GameLocation
 * @property {string} id
 * @property {string} game_id
 * @property {string} name
 * @property {string} description
 * @property {string} created_at
 * @property {string} [updated_at]
 * @property {string} [deleted_at]
 */

/**
 * @typedef {Object} GameLocationLink
 * @property {string} id
 * @property {string} game_id
 * @property {string} from_adventure_game_location_id
 * @property {string} to_adventure_game_location_id
 * @property {string} name
 * @property {string} description
 * @property {string} created_at
 * @property {string} [updated_at]
 * @property {string} [deleted_at]
 */

/**
 * @typedef {Object} GameCreatureInstance
 * @property {string} id
 * @property {string} game_id
 * @property {string} game_creature_id
 * @property {string} game_instance_id
 * @property {string} game_location_instance_id
 * @property {boolean} is_alive
 * @property {string} created_at
 * @property {string} [updated_at]
 * @property {string} [deleted_at]
 */

/**
 * @typedef {Object} GameItemInstance
 * @property {string} id
 * @property {string} game_id
 * @property {string} game_item_id
 * @property {string|null} location_id
 * @property {string|null} character_id
 * @property {string|null} creature_id
 * @property {boolean} is_equipped
 * @property {boolean} is_used
 * @property {number|null} uses_remaining
 * @property {string} created_at
 * @property {string} [updated_at]
 * @property {string} [deleted_at]
 */

/**
 * @typedef {Object} GameLocationLinkRequirement
 * @property {string} id
 * @property {string} game_id
 * @property {string} game_location_link_id
 * @property {string} game_item_id
 * @property {number} quantity
 * @property {string} created_at
 * @property {string} [updated_at]
 * @property {string} [deleted_at]
 */ 