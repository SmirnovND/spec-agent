package cli

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/SmirnovND/spec-agent/internal/config"
	"github.com/SmirnovND/spec-agent/internal/spec"
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("port", "p", "8080", "порт для сервера")
	serveCmd.Flags().String("host", "localhost", "хост для привязки")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Запустить встроенный веб-сервер для просмотра спецификаций",
	Long: `
Команда serve:
- генерирует HTML спецификации (если нет в .spec_agent/build/)
- запускает встроенный HTTP сервер
- обслуживает файлы из .spec_agent/build/
- доступна по http://localhost:8080
- прекращает работу по Ctrl+C
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetString("port")
		host, _ := cmd.Flags().GetString("host")

		buildDir := filepath.Join(".spec_agent", "build")
		indexPath := filepath.Join(buildDir, "index.html")

		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			fmt.Println("📝 Спецификации не сгенерированы. Генерирую...")
			fmt.Println()

			if err := generateSpecs(); err != nil {
				return fmt.Errorf("ошибка при генерации: %w", err)
			}

			fmt.Println()
		}

		return serveFiles(host, port, buildDir)
	},
}

func generateSpecs() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("не удалось загрузить config.yaml: %w", err)
	}

	if len(cfg.Roots) == 0 {
		return fmt.Errorf("в config.yaml не указаны roots")
	}

	specFiles, err := findSpecsNearRoots(cfg.Roots)
	if err != nil {
		return err
	}

	if len(specFiles) == 0 {
		return fmt.Errorf("не найдено ни одной спецификации рядом с roots")
	}

	referenced := spec.CollectAllReferences(specFiles)
	rootSpecs := spec.FindRootSpecs(specFiles, referenced)
	if len(rootSpecs) == 0 {
		return fmt.Errorf("не удалось определить корневые спецификации")
	}

	fmt.Printf("🌳 Найдено %d корневых спецификаций\n", len(rootSpecs))

	graph, err := spec.BuildGraphFromRoots(rootSpecs)
	if err != nil {
		return err
	}

	fmt.Printf("📊 Граф содержит %d узлов и %d ребер\n", len(graph.Nodes), len(graph.Edges))

	buildDir := filepath.Join(".spec_agent", "build")
	if err := spec.ExportToHTML(graph, buildDir); err != nil {
		return fmt.Errorf("ошибка при экспорте: %w", err)
	}

	fmt.Println("✅ HTML сгенерирован успешно!")
	return nil
}

func serveFiles(host, port, buildDir string) error {
	absPath, _ := filepath.Abs(buildDir)

	fs := http.FileServer(http.Dir(buildDir))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			r.URL.Path = "/index.html"
		}
		fs.ServeHTTP(w, r)
	})

	addr := net.JoinHostPort(host, port)
	server := &http.Server{
		Addr: addr,
	}

	fmt.Println()
	fmt.Printf("🚀 Веб-сервер запущен!\n")
	fmt.Printf("🌐 Откройте в браузере: http://%s:%s\n", host, port)
	fmt.Printf("📂 Обслуживаются файлы из: %s\n", absPath)
	fmt.Println()
	fmt.Println("Нажмите Ctrl+C для выключения сервера")
	fmt.Println()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println()
		fmt.Println("⏹️  Сервер выключен")
		os.Exit(0)
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("ошибка сервера: %w", err)
	}

	return nil
}
