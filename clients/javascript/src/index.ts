export { browser, Browser, LaunchOptions } from './browser';
export { Page } from './page';
export { BrowserContext } from './context';
export { Vibe, FindOptions } from './vibe';
export { Element, BoundingBox, ElementInfo, ActionOptions, SelectorOptions, FluentElement, fluent } from './element';
export { ElementList, FilterOptions } from './element-list';

// Sync API
export { browserSync, VibeSync, ElementSync } from './sync';

// Error types
export {
  ConnectionError,
  TimeoutError,
  ElementNotFoundError,
  BrowserCrashedError,
} from './utils/errors';
