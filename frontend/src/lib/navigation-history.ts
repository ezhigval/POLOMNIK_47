export const NAV_PREV_PATH_KEY = "palomnik_prev_path";
export const NAV_CURRENT_PATH_KEY = "palomnik_current_path";

export function readPrevPath(): string | null {
  if (typeof window === "undefined") {
    return null;
  }
  return sessionStorage.getItem(NAV_PREV_PATH_KEY);
}

/** Whether browser back is likely to stay on-site and leave the current page. */
export function canUseHistoryBack(currentPath: string): boolean {
  if (typeof window === "undefined") {
    return false;
  }
  if (window.history.length <= 1) {
    return false;
  }

  const prevPath = readPrevPath();
  if (prevPath && prevPath !== currentPath) {
    return true;
  }

  const referrer = document.referrer;
  if (!referrer) {
    return false;
  }

  try {
    const ref = new URL(referrer);
    if (ref.origin !== window.location.origin) {
      return false;
    }
    return ref.pathname !== currentPath;
  } catch {
    return false;
  }
}
