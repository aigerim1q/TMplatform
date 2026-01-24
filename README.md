# ЖЦП Parser - AI-Powered Document Parser

> Автоматический парсер документов жизненного цикла проекта (ЖЦП) с использованием AI и автоматическим назначением ответственных лиц.

## 🎯 Описание

Парсер автоматически извлекает структурированную информацию из PDF и DOCX документов:

- Фазы проекта
- Задачи и подзадачи
- Временные рамки (даты начала/окончания)
- **Автоматическое назначение ответственных** на основе анализа задач
- Зависимости между задачами

Результаты сохраняются в JSON и SQLite базу данных.

## ✨ Ключевые возможности

- 📄 **Парсинг PDF/DOCX** - извлечение текста из документов
- 🤖 **AI-анализ** - использование LLM для структурирования данных
- 👥 **Автоназначение ответственных** - 14 специалистов с разными должностями
- 🔄 **Fallback система** - поддержка OpenAI, Anthropic, Ollama, DeepSeek
- ✅ **Валидация данных** - проверка структуры и качества
- 💾 **Двойное сохранение** - JSON файлы + SQLite БД

## 🚀 Быстрый старт

### 1. Требования

- Go 1.23+ (или 1.24+)
- Git

### 2. Установка

```bash
git clone <your-repo-url>
cd parsing-ai/zhcp-parser-go
go mod tidy
```

### 3. Конфигурация LLM

Отредактируйте `configs/llm_config.yaml`:

```yaml
providers:
  ollama:
    enabled: true # Локальный LLM (по умолчанию)
    model: "llama3"
    base_url: "http://localhost:11434"
```

Для облачных провайдеров:

```bash
export OPENAI_API_KEY="your-key"
export ANTHROPIC_API_KEY="your-key"
```

### 4. Запуск

```bash
# Парсинг документа
go run cmd/zhcp-parser/main.go parse path/to/document.pdf

# С опциями
go run cmd/zhcp-parser/main.go parse document.pdf --validate --enrich --output result.json

# Тест на примере
go run cmd/zhcp-parser/main.go parse testdata/sample_project.txt
```

### 5. Компиляция

```bash
go build -o zhcp-parser cmd/zhcp-parser/main.go
./zhcp-parser parse document.pdf
```

## 👥 Автоматическое назначение ответственных

Система включает 14 вымышленных сотрудников:

- 👔 Руководитель проекта
- 📊 Бизнес-аналитик
- 🏗️ Архитектор решений
- ⚙️ Backend разработчик
- 🎨 Frontend разработчик
- 🔄 Fullstack разработчик
- 🧪 Тестировщик
- 🚀 DevOps инженер
- 🎨 UI/UX дизайнер
- 🤖 **AI интегратор** ⭐
- 📈 Data Scientist
- 📝 Технический писатель
- 🔒 Специалист по безопасности
- 📱 Мобильный разработчик

LLM автоматически анализирует содержание каждой задачи и назначает подходящего специалиста на основе ключевых слов.

### Примеры назначения:

- "Разработка REST API" → Backend разработчик
- "Интеграция ChatGPT" → AI интегратор
- "Дизайн интерфейса" → UI/UX дизайнер
- "Настройка CI/CD" → DevOps инженер

Настройка пула: `prompts/employee_pool.json`

## 📁 Структура проекта

```
zhcp-parser-go/
├── cmd/zhcp-parser/          # CLI приложение
├── configs/                  # Конфигурации LLM
├── internal/
│   ├── ai/                   # AI интеграция и промпты
│   ├── parser/               # Парсеры документов
│   ├── parsers/              # PDF/DOCX экстракторы
│   ├── transformers/         # Обработка данных
│   ├── validators/           # Валидация
│   └── storage/              # SQLite хранилище
├── prompts/                  # Промпты и пул сотрудников
├── testdata/                 # Тестовые документы
├── QUICKSTART.md             # Быстрый старт
└── README.md                 # Документация
```

## 📊 Пример результата

```json
{
  "project": {
    "title": "Разработка веб-платформы",
    "phases": [
      {
        "name": "Разработка Backend",
        "tasks": [
          {
            "name": "Разработка REST API",
            "responsible_persons": [
              {
                "name": "Иван Волков",
                "role": "Backend разработчик"
              }
            ]
          },
          {
            "name": "Интеграция AI для рекомендаций",
            "responsible_persons": [
              {
                "name": "Роман Белов",
                "role": "AI интегратор"
              }
            ]
          }
        ]
      }
    ]
  }
}
```

## 🔧 Поддерживаемые LLM провайдеры

- **Ollama** (локальный) - Llama 3, Mistral и др.
- **OpenAI** - GPT-4, GPT-4 Turbo
- **Anthropic** - Claude 3 (Sonnet, Opus)
- **DeepSeek** - DeepSeek Chat

Система автоматически переключается между провайдерами при сбоях.

## 📖 Документация

- **[QUICKSTART.md](zhcp-parser-go/QUICKSTART.md)** - Быстрый старт
- **[README.md](zhcp-parser-go/README.md)** - Подробная документация
- **[configs/llm_config.yaml](zhcp-parser-go/configs/llm_config.yaml)** - Пример конфигурации

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch
3. Commit изменений
4. Push в branch
5. Создайте Pull Request

## 📝 Лицензия

MIT License

## 🆘 Поддержка

При возникновении проблем создайте Issue в GitHub репозитории.

---

**Разработано с использованием:** Go, Cobra CLI, SQLite, OpenAI API, Anthropic API, Ollama
