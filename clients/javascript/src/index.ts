export { browser, Browser, LaunchOptions } from './browser';
export { Page } from './page';
export { BrowserContext } from './context';
export { Vibe, FindOptions } from './vibe';
export { Element, BoundingBox, ElementInfo, ActionOptions } from './element';

// Sync API
export { browserSync, VibeSync, ElementSync } from './sync';

// Error types
export {
  ConnectionError,
  TimeoutError,
  ElementNotFoundError,
  BrowserCrashedError,
} from './utils/errors';
