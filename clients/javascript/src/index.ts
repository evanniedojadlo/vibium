export { browser, Browser, LaunchOptions } from './browser';
export { Page, Keyboard, Mouse, Touch, ScreenshotOptions, A11yNode } from './page';
export { BrowserContext, Cookie, SetCookieParam, StorageState, OriginState } from './context';
export { Vibe, FindOptions } from './vibe';
export { Element, BoundingBox, ElementInfo, ActionOptions, SelectorOptions, FluentElement, fluent } from './element';
export { ElementList, FilterOptions } from './element-list';
export { Route } from './route';
export { Request, Response } from './network';
export { Dialog } from './dialog';
export { ConsoleMessage } from './console';
export { Download } from './download';
export { WebSocketInfo } from './websocket';

// Sync API
export { browserSync, VibeSync, ElementSync } from './sync';

// Error types
export {
  ConnectionError,
  TimeoutError,
  ElementNotFoundError,
  BrowserCrashedError,
} from './utils/errors';
