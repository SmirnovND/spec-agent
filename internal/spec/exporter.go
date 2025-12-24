package spec

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
)

func ExportToHTML(graph *Graph, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию: %w", err)
	}

	specs := make(map[string]*Spec)
	for _, node := range graph.Nodes {
		spec, err := ParseFile(node.Path)
		if err != nil {
			continue
		}
		specs[node.Path] = spec
	}

	indexHTML := generateIndexHTML(specs, graph)
	if err := os.WriteFile(filepath.Join(outputDir, "index.html"), []byte(indexHTML), 0644); err != nil {
		return fmt.Errorf("не удалось записать index.html: %w", err)
	}

	for path, spec := range specs {
		filename := generateFilename(path)
		specHTML := generateSpecHTML(spec, specs)
		if err := os.WriteFile(filepath.Join(outputDir, filename), []byte(specHTML), 0644); err != nil {
			return fmt.Errorf("не удалось записать %s: %w", filename, err)
		}
	}

	return nil
}

func generateIndexHTML(specs map[string]*Spec, graph *Graph) string {
	toc := ""
	for path, spec := range specs {
		if spec.Title == "" {
			spec.Title = filepath.Base(path)
		}
		filename := generateFilename(path)
		toc += fmt.Sprintf(`    <li><a href="%s">%s</a></li>
`, filename, html.EscapeString(spec.Title))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Спецификации архитектуры</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', sans-serif;
            line-height: 1.6;
            color: #333;
            background: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 40px 20px;
        }
        header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 40px 20px;
            border-radius: 8px;
            margin-bottom: 40px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
        }
        header p {
            font-size: 1.1em;
            opacity: 0.95;
        }
        .main-content {
            display: grid;
            grid-template-columns: 250px 1fr;
            gap: 40px;
        }
        .sidebar {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
            height: fit-content;
            position: sticky;
            top: 20px;
        }
        .sidebar h2 {
            font-size: 1.2em;
            margin-bottom: 20px;
            color: #667eea;
            border-bottom: 2px solid #667eea;
            padding-bottom: 10px;
        }
        .sidebar ul {
            list-style: none;
        }
        .sidebar li {
            margin-bottom: 10px;
        }
        .sidebar a {
            color: #667eea;
            text-decoration: none;
            padding: 8px 12px;
            display: block;
            border-radius: 4px;
            transition: all 0.2s;
        }
        .sidebar a:hover {
            background: #f0f0f0;
            transform: translateX(4px);
        }
        .content {
            background: white;
            padding: 40px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        }
        .welcome {
            text-align: center;
            padding: 60px 20px;
        }
        .welcome h2 {
            font-size: 2em;
            color: #667eea;
            margin-bottom: 20px;
        }
        .welcome p {
            font-size: 1.1em;
            color: #666;
        }
        @media (max-width: 768px) {
            .main-content {
                grid-template-columns: 1fr;
            }
            .sidebar {
                position: static;
            }
            header h1 {
                font-size: 1.8em;
            }
        }
    </style>
</head>
<body>
    <header>
        <div class="container">
            <h1>📚 Спецификации архитектуры</h1>
            <p>Документация компонентов и их взаимодействия</p>
        </div>
    </header>
    <div class="container">
        <div class="main-content">
            <div class="sidebar">
                <h2>Спецификации</h2>
                <ul>
%s
                </ul>
            </div>
            <div class="content">
                <div class="welcome">
                    <h2>Добро пожаловать!</h2>
                    <p>Выберите спецификацию из меню слева для просмотра деталей</p>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`, toc)
}

func generateSpecHTML(spec *Spec, allSpecs map[string]*Spec) string {
	title := spec.Title
	if title == "" {
		title = filepath.Base(spec.Path)
	}

	contentHTML := markdownToHTML(spec.Content)

	navigation := ""
	for _, link := range spec.Links {
		filename := generateFilenameFromRelative(spec.Path, link.Path)
		navigation += fmt.Sprintf(`        <li><a href="%s">%s</a></li>
`, filename, html.EscapeString(link.Title))
	}

	navSection := ""
	if navigation != "" {
		navSection = fmt.Sprintf(`        <div class="navigation">
            <h3>Связанные спецификации</h3>
            <ul>
%s
            </ul>
        </div>
`, navigation)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - Спецификации</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', sans-serif;
            line-height: 1.8;
            color: #333;
            background: #f5f5f5;
        }
        .container {
            max-width: 900px;
            margin: 0 auto;
            padding: 40px 20px;
        }
        header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 30px 20px;
            margin-bottom: 40px;
            border-radius: 8px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        header .back-link {
            display: inline-block;
            margin-bottom: 15px;
            color: rgba(255,255,255,0.9);
            text-decoration: none;
            font-size: 0.95em;
            transition: color 0.2s;
        }
        header .back-link:hover {
            color: white;
        }
        header h1 {
            font-size: 2em;
            margin-bottom: 5px;
        }
        .content {
            background: white;
            padding: 40px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        }
        .content h1 { font-size: 2em; margin-top: 30px; margin-bottom: 20px; color: #333; border-bottom: 2px solid #667eea; padding-bottom: 10px; }
        .content h1:first-child { margin-top: 0; }
        .content h2 { font-size: 1.5em; margin-top: 25px; margin-bottom: 15px; color: #555; }
        .content h3 { font-size: 1.2em; margin-top: 20px; margin-bottom: 12px; color: #666; }
        .content p { margin-bottom: 15px; }
        .content ul, .content ol { margin-left: 30px; margin-bottom: 15px; }
        .content li { margin-bottom: 8px; }
        .content code { background: #f4f4f4; padding: 2px 6px; border-radius: 3px; font-family: 'Monaco', 'Menlo', monospace; font-size: 0.95em; color: #d73a49; }
        .content pre { background: #f4f4f4; padding: 15px; border-radius: 5px; overflow-x: auto; margin-bottom: 15px; }
        .content pre code { color: #333; padding: 0; background: none; }
        .content blockquote { border-left: 4px solid #667eea; padding-left: 15px; margin-left: 0; margin-bottom: 15px; color: #666; font-style: italic; }
        .content table { width: 100%%; border-collapse: collapse; margin-bottom: 15px; }
        .content table th, .content table td { border: 1px solid #ddd; padding: 12px; text-align: left; }
        .content table th { background: #f9f9f9; font-weight: 600; }
        .navigation {
            margin-top: 40px;
            padding-top: 30px;
            border-top: 2px solid #f0f0f0;
        }
        .navigation h3 {
            color: #667eea;
            margin-bottom: 15px;
        }
        .navigation ul {
            list-style: none;
            margin-left: 0;
        }
        .navigation li {
            margin-bottom: 10px;
        }
        .navigation a {
            color: #667eea;
            text-decoration: none;
            padding: 8px 12px;
            display: inline-block;
            border-radius: 4px;
            transition: all 0.2s;
            border: 1px solid #667eea;
        }
        .navigation a:hover {
            background: #667eea;
            color: white;
        }
        @media (max-width: 768px) {
            .container { padding: 20px 15px; }
            header h1 { font-size: 1.5em; }
            .content { padding: 20px; }
        }
    </style>
</head>
<body>
    <header>
        <div class="container">
            <a href="index.html" class="back-link">← Вернуться к спецификациям</a>
            <h1>%s</h1>
        </div>
    </header>
    <div class="container">
        <div class="content">
            %s
%s
        </div>
    </div>
</body>
</html>`, html.EscapeString(title), html.EscapeString(title), contentHTML, navSection)
}

func generateFilename(path string) string {
	return filepath.Base(path) + ".html"
}

func generateFilenameFromRelative(basePath, relativePath string) string {
	dir := filepath.Dir(basePath)
	absPath := filepath.Join(dir, relativePath)
	return filepath.Base(absPath) + ".html"
}

func markdownToHTML(content string) string {
	lines := strings.Split(content, "\n")
	var result strings.Builder
	var inCode bool
	var inList bool

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			if inCode {
				result.WriteString("</pre>\n")
				inCode = false
			} else {
				result.WriteString("<pre><code>")
				inCode = true
			}
			continue
		}

		if inCode {
			result.WriteString(html.EscapeString(line))
			result.WriteString("\n")
			continue
		}

		if trimmed == "" {
			if inList {
				result.WriteString("</ul>\n")
				inList = false
			}
			result.WriteString("\n")
			continue
		}

		if strings.HasPrefix(trimmed, "# ") {
			result.WriteString(fmt.Sprintf("<h1>%s</h1>\n", html.EscapeString(strings.TrimPrefix(trimmed, "# "))))
		} else if strings.HasPrefix(trimmed, "## ") {
			result.WriteString(fmt.Sprintf("<h2>%s</h2>\n", html.EscapeString(strings.TrimPrefix(trimmed, "## "))))
		} else if strings.HasPrefix(trimmed, "### ") {
			result.WriteString(fmt.Sprintf("<h3>%s</h3>\n", html.EscapeString(strings.TrimPrefix(trimmed, "### "))))
		} else if strings.HasPrefix(trimmed, "- ") {
			if !inList {
				result.WriteString("<ul>\n")
				inList = true
			}
			item := strings.TrimPrefix(trimmed, "- ")
			result.WriteString(fmt.Sprintf("<li>%s</li>\n", html.EscapeString(item)))
		} else if strings.HasPrefix(trimmed, "> ") {
			result.WriteString(fmt.Sprintf("<blockquote>%s</blockquote>\n", html.EscapeString(strings.TrimPrefix(trimmed, "> "))))
		} else {
			if inList {
				result.WriteString("</ul>\n")
				inList = false
			}
			result.WriteString(fmt.Sprintf("<p>%s</p>\n", html.EscapeString(trimmed)))
		}
	}

	if inList {
		result.WriteString("</ul>\n")
	}
	if inCode {
		result.WriteString("</code></pre>\n")
	}

	return result.String()
}
