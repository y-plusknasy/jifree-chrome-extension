export interface AnalyzeRequest {
  html: string;
  selection: string;
  prefix: string;
  suffix: string;
  auth: {
    shared_secret: string;
    user_id: string;
  };
}

export interface AnalyzeResponse {
  tokens: Token[];
}

export interface Token {
  base_form: string;
  reading: string;
  part_of_speech: string;
}

// 拡張機能内部のメッセージング用
export interface MessageRequest {
  type: 'ANALYZE_TEXT';
  payload: {
    html: string;
    selection: string;
    prefix: string;
    suffix: string;
  };
}

export interface MessageResponse {
  success: boolean;
  data?: AnalyzeResponse;
  error?: string;
}
