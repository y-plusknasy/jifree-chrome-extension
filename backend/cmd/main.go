package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/y-plusknasy/jifree-chrome-extension/backend" // init()を呼び出すためにインポート
)

func main() {
	// ローカル開発用サーバーの起動
	// FUNCTION_TARGET環境変数で指定された関数を実行する
	// docker-compose.ymlで FUNCTION_TARGET=AnalyzeText と指定されている想定
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
