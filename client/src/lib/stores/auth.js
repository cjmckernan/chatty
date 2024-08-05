import { writable } from 'svelte/store';

function isLocalStorageAvailable() {
  try {
    const testKey = '__test__';
    localStorage.setItem(testKey, testKey);
    localStorage.removeItem(testKey);
    return true;
  } catch (e) {
    return false;
  }
}

function getSessionIdFromLocalStorage() {
  if (isLocalStorageAvailable()) {
    return localStorage.getItem('sessionId') || null;
  } else {
    return null;
  }
}

export const sessionId = writable(getSessionIdFromLocalStorage());

sessionId.subscribe(value => {
  if (isLocalStorageAvailable()) {
    if (value) {
      localStorage.setItem('sessionId', value);
    } else {
      localStorage.removeItem('sessionId');
    }
  }
});

export function setSessionId(xSessionId) {
  sessionId.set(xSessionId);
}

export function destroySessionId() {
  sessionId.set(null);
}

