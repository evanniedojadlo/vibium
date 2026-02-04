// BiDi Protocol Types

export interface BiDiCommand {
  id: number;
  method: string;
  params: Record<string, unknown>;
}

export interface BiDiResponse {
  id: number;
  type: 'success' | 'error';
  result?: unknown;
  error?: string | { error?: string; message?: string }; // Can be flat (spec) or nested (current impl)
  message?: string;
}

export interface BiDiEvent {
  method: string;
  params: Record<string, unknown>;
}

export interface BiDiError {
  error: string;
  message: string;
  stacktrace?: string;
}

export type BiDiMessage = BiDiResponse | BiDiEvent;

export function isResponse(msg: BiDiMessage): msg is BiDiResponse {
  return 'id' in msg;
}

export function isEvent(msg: BiDiMessage): msg is BiDiEvent {
  return !('id' in msg) && 'method' in msg;
}

// Session types
export interface SessionStatus {
  ready: boolean;
  message: string;
}

// Browsing Context types
export interface BrowsingContextInfo {
  context: string;
  url: string;
  children: BrowsingContextInfo[];
  parent?: string;
}

export interface BrowsingContextTree {
  contexts: BrowsingContextInfo[];
}

export interface NavigationResult {
  navigation: string;
  url: string;
}

export interface ScreenshotResult {
  data: string; // base64 PNG
}
