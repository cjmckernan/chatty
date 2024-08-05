import { writable } from 'svelte/store';

export const username = writable(null)

export function setUsername(requestUsername) {
  username.set(requestUsername);
}
