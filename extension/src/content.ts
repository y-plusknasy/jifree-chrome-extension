import { MessageRequest, MessageResponse, AnalyzeResponse, Token } from './types';

let currentPopup: HTMLElement | null = null;

document.addEventListener('mouseup', async () => {
  const selection = window.getSelection();
  if (!selection || selection.isCollapsed) {
    removePopup();
    return;
  }

  const text = selection.toString().trim();
  if (!text) return;

  // 選択範囲の情報を取得
  const range = selection.getRangeAt(0);
  const rect = range.getBoundingClientRect();

  // 文脈情報の抽出
  const { prefix, suffix, html } = extractContext(range);

  // APIリクエスト送信
  try {
    const message: MessageRequest = {
      type: 'ANALYZE_TEXT',
      payload: {
        html,
        selection: text,
        prefix,
        suffix
      }
    };

    const response = await chrome.runtime.sendMessage(message) as MessageResponse;
    
    if (response.success && response.data) {
      showPopup(rect, response.data.tokens);
    } else {
      console.error('Analysis failed:', response.error);
    }
  } catch (error) {
    console.error('Communication error:', error);
  }
});

// ポップアップ以外の場所をクリックしたら閉じる
document.addEventListener('mousedown', (e) => {
  if (currentPopup && !currentPopup.contains(e.target as Node)) {
    removePopup();
  }
});

function removePopup() {
  if (currentPopup) {
    currentPopup.remove();
    currentPopup = null;
  }
}

function extractContext(range: Range): { prefix: string, suffix: string, html: string } {
  // 親要素のHTMLを取得
  // startContainerがテキストノードの場合はその親要素、要素ノードの場合はその要素を使用
  let container = range.commonAncestorContainer;
  if (container.nodeType === Node.TEXT_NODE) {
    container = container.parentElement as HTMLElement;
  }

  // コンテキストが短すぎる場合、親要素へ遡る
  // sectionタグ、またはbodyタグに到達するまで、かつ文字数が50文字未満の場合に遡る
  while (
    container.textContent && 
    container.textContent.length < 50 && 
    container.parentElement && 
    (container as HTMLElement).tagName.toLowerCase() !== 'section' &&
    (container as HTMLElement).tagName.toLowerCase() !== 'body'
  ) {
    container = container.parentElement;
  }

  const htmlContent = (container as HTMLElement).outerHTML || '';

  // Prefix抽出 (簡易版: コンテナ内のテキストから切り出し)
  // 本来はRangeを使って厳密に取るべきだが、まずはtextContentベースで実装
  const fullText = container.textContent || '';
  const selectionText = range.toString();
  
  // 単純なindexOfだと同じ単語が複数ある場合にずれる可能性があるが、
  // バックエンド側でもPrefix/Suffixマッチングを行うため、ある程度の周辺テキストが取れれば良い
  // ここでは簡易的に、選択範囲の前方10文字、後方10文字を取得する
  
  // より正確に取るためにRangeを拡張する
  const preRange = range.cloneRange();
  preRange.collapse(true); // 先頭に潰す
  preRange.setStart(container, 0); // コンテナの先頭から
  const prefixFull = preRange.toString();
  const prefix = prefixFull.slice(-10); // 直前10文字

  const postRange = range.cloneRange();
  postRange.collapse(false); // 末尾に潰す
  postRange.setEndAfter(container); // コンテナの末尾まで
  const suffixFull = postRange.toString();
  const suffix = suffixFull.slice(0, 10); // 直後10文字

  return {
    html: htmlContent,
    prefix,
    suffix
  };
}

function showPopup(rect: DOMRect, tokens: Token[]) {
  removePopup();

  // スタイルの注入（初回のみ）
  if (!document.getElementById('jifree-style')) {
    const style = document.createElement('style');
    style.id = 'jifree-style';
    style.textContent = `
      .jifree-popup {
        position: absolute;
        background-color: white;
        border: 2px solid #333;
        border-radius: 30px;
        padding: 16px 24px;
        box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        z-index: 999999;
        max-width: 400px;
        font-family: sans-serif;
        font-size: 20px;
        line-height: 1.6;
        color: #333;
      }
      .jifree-popup::before {
        content: '';
        position: absolute;
        top: -14px;
        left: 30px;
        border: 7px solid transparent;
        border-bottom-color: #333;
      }
      .jifree-popup::after {
        content: '';
        position: absolute;
        top: -10px;
        left: 32px;
        border: 5px solid transparent;
        border-bottom-color: white;
      }
      .jifree-token {
        display: inline-block;
        margin-right: 6px;
        text-align: center;
        color: #333;
      }
      .jifree-token ruby {
        font-weight: bold;
        color: #333;
      }
      .jifree-token rt {
        font-size: 0.75em;
        color: #333;
        font-weight: bold;
      }
    `;
    document.head.appendChild(style);
  }

  const popup = document.createElement('div');
  popup.className = 'jifree-popup';
  popup.style.left = `${window.scrollX + rect.left}px`;
  popup.style.top = `${window.scrollY + rect.bottom + 15}px`; // 選択範囲の下に表示（吹き出しの分調整）

  // トークンを表示
  tokens.forEach(token => {
    const span = document.createElement('span');
    span.className = 'jifree-token';

    const ruby = document.createElement('ruby');
    ruby.textContent = token.base_form;
    
    const rt = document.createElement('rt');
    rt.textContent = token.reading;

    ruby.appendChild(rt);
    span.appendChild(ruby);
    popup.appendChild(span);
  });

  document.body.appendChild(popup);
  currentPopup = popup;
}
