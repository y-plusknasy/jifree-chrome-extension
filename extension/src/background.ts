import { CONFIG } from './utils/config';
import { MessageRequest, MessageResponse, AnalyzeRequest, AnalyzeResponse } from './types';

// UUID生成 (簡易版)
function generateUUID(): string {
  return crypto.randomUUID();
}

// ユーザーIDの取得または生成
async function getUserId(): Promise<string> {
  const key = 'jifree_user_id';
  const result = await chrome.storage.local.get(key);
  if (result[key]) {
    return result[key] as string;
  }

  const newId = generateUUID();
  await chrome.storage.local.set({ [key]: newId });
  return newId;
}

// メッセージリスナー
chrome.runtime.onMessage.addListener((
  message: MessageRequest,
  sender,
  sendResponse: (response: MessageResponse) => void
) => {
  if (message.type === 'ANALYZE_TEXT') {
    handleAnalyzeText(message)
      .then(response => sendResponse({ success: true, data: response }))
      .catch(error => sendResponse({ success: false, error: error.message }));
    
    return true; // 非同期レスポンスのためにtrueを返す
  }
});

async function handleAnalyzeText(message: MessageRequest): Promise<AnalyzeResponse> {
  const userId = await getUserId();
  
  const requestBody: AnalyzeRequest = {
    html: message.payload.html,
    selection: message.payload.selection,
    prefix: message.payload.prefix,
    suffix: message.payload.suffix,
    auth: {
      shared_secret: CONFIG.SHARED_SECRET,
      user_id: userId
    }
  };

  const response = await fetch(CONFIG.API_ENDPOINT, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(requestBody)
  });

  if (!response.ok) {
    if (response.status === 429) {
      throw new Error('リクエストが多すぎます。少し待ってから再試行してください。');
    }
    throw new Error(`API Error: ${response.statusText}`);
  }

  return await response.json();
}
