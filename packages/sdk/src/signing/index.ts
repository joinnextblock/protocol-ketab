import type { BaseNostrEvent } from "@ketab/core";
import { finalizeEvent, getPublicKey as getPubkey, verifyEvent as verify } from "nostr-tools/pure";
import { bytesToHex, hexToBytes } from "@noble/hashes/utils";

/**
 * Sign a Nostr event with a secret key
 * @param event - Unsigned event (without id and sig)
 * @param secret_key - Secret key as Uint8Array (32 bytes)
 */
export function sign_event(
  event: Omit<BaseNostrEvent, "id" | "sig">,
  secret_key: Uint8Array
): BaseNostrEvent {
  return finalizeEvent(event, secret_key) as BaseNostrEvent;
}

/**
 * Get public key (hex) from secret key
 * @param secret_key - Secret key as Uint8Array (32 bytes)
 */
export function get_public_key(secret_key: Uint8Array): string {
  return getPubkey(secret_key);
}

/**
 * Verify an event signature
 */
export function verify_event(event: BaseNostrEvent): boolean {
  return verify(event);
}

/**
 * Convert hex string to Uint8Array
 */
export function hex_to_secret_key(hex: string): Uint8Array {
  return hexToBytes(hex);
}

/**
 * Convert Uint8Array to hex string
 */
export function secret_key_to_hex(secret_key: Uint8Array): string {
  return bytesToHex(secret_key);
}
