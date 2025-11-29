// 本番環境ではビルド時に置換するか、環境変数から注入することを推奨
// 今回はスモールスタートのためハードコード（ただしリポジトリ公開時は注意）
export const CONFIG = {
  API_ENDPOINT: 'http://localhost:8080', // ローカル開発用
  // API_ENDPOINT: 'https://YOUR_REGION-YOUR_PROJECT.cloudfunctions.net/AnalyzeText', // 本番用
  SHARED_SECRET: 'test-secret', // バックエンドの環境変数と一致させる
};
