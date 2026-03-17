# AGENTS.md

Este proyecto usa **OpenSpec** para SDD (Spec-Driven Development) y varios skills para diferentes tareas.

## Tabla de Skills Disponibles

| Skill | Descripción | Trigger | Scope | Path |
|--------|-------------|---------|-------|------|
| **openspec-propose** | Crear nuevo cambio con specs, diseño y tareas | "quiero agregar algo nuevo", "necesito implementar X", "nueva feature" | Crear specs y tareas para nueva funcionalidad | [.opencode/skills/openspec-propose/SKILL.md](.opencode/skills/openspec-propose/SKILL.md) |
| **openspec-apply-change** | Implementar tareas de un cambio existente | "implementa el cambio", "haz las tareas", "/opsx-apply" | Ejecutar tareas de un change en openspec | [.opencode/skills/openspec-apply-change/SKILL.md](.opencode/skills/openspec-apply-change/SKILL.md) |
| **openspec-explore** | Explorar ideas antes de crear cambio | "cómo podría hacer X", "quiero pensar en", "explora la arquitectura de" | Análisis previo a proposal | [.opencode/skills/openspec-explore/SKILL.md](.opencode/skills/openspec-explore/SKILL.md) |
| **openspec-archive-change** | Archivar cambio completado | "archiva el cambio", "change已完成" | Finalizar y documentar cambio | [.opencode/skills/openspec-archive-change/SKILL.md](.opencode/skills/openspec-archive-change/SKILL.md) |
| **golang-pro** | Patrones Go avanzados, concurrencia, gRPC | "patrones concurrentes", "goroutines", "microservicio", "optimizar Go" | Código Go avanzado | [.agents/skills/golang-pro/SKILL.md](.agents/skills/golang-pro/SKILL.md) |
| **mermaid-diagrams** | Crear diagramas profesionales | "diagrama de", "visualiza", "mermaid", "arquitectura" | Diagramas para documentación | [.agents/skills/mermaid-diagrams/SKILL.md](.agents/skills/mermaid-diagrams/SKILL.md) |
| **find-skills** | Descubrir skills disponibles | "hay skill para X", "cómo hago para", "existe skill" | Buscar skills instalados | (skill del sistema) |

## Cuándo Invocar Cada Skill

### openspec-propose

Usar cuando:
- Se quiere agregar una nueva capability al sistema
- Se necesita crear specs para algo que no existe
- Se va a implementar algo desde cero

```bash
# Desde terminal
openspec new change "<nombre-del-cambio>"

# O pedir al agente
openspec-propose
```

### openspec-apply-change

Usar cuando:
- Ya existe un cambio con tareas pendientes
- Se quiere implementar el contenido de un change

```bash
# Desde terminal
openspec apply --change "<nombre>"

# O pedir al agente
openspec-apply-change
```

### openspec-explore

Usar cuando:
- Se quiere analizar un problema antes de actuar
- Se necesita aclarar requisitos
- Se va a hacer refactoring grande

```bash
openspec-explore
```

### openspec-archive-change

Usar cuando:
- Un cambio está completo
- Se quiere documentar que ya se terminó

```bash
openspec-archive-change
```

### golang-pro

Usar cuando:
- Se necesita código Go con patrones avanzados
- Concurrencia, goroutines, channels
- gRPC o REST APIs
- CLI tools
- Testing avanzado

### mermaid-diagrams

Usar cuando:
- Se necesita documentar arquitectura
- Diagramas de secuencia para flujos
- Diagramas de clase para modelos
- ERDs para base de datos
- Flowcharts para procesos

## OpenSpec - SDD del Proyecto

### Estructura

```
openspec/
├── specs/                    # Specs implementadas (fuente de verdad)
│   ├── domain-models/
│   ├── pipeline-engine/
│   ├── pipeline-stages/
│   ├── webhook-handler/
│   ├── whatsapp-client/
│   └── observability/
│
├── changes/                  # Cambios con tareas pendientes
│   └── implement-tts-pipeline/
│       ├── proposal.md
│       ├── design.md
│       ├── specs.md
│       └── tasks.md
│
└── config.yaml              # Config de OpenSpec
```

### Flujo de Trabajo

1. **Crear cambio**: `openspec new change "<nombre>"`
2. **Implementar**: `openspec apply --change "<nombre>"`
3. **Archivar**: `openspec archive --change "<nombre>"`

### Specs

- `openspec/specs/` - Specs canonicales del sistema (fuente de verdad)
- `docs/` - Documentación de referencia (Postman, prompts)

**Importante**: `openspec/specs/` es la fuente de verdad. Los otros directorios son referencia histórica.

## Reglas

1. **Siempre** crear specs en `openspec/specs/` para nueva funcionalidad
2. **Siempre** usar changes para trabajo en progreso
3. **Siempre** actualizar tasks al completar items
4. **Nunca** modificar specs sin pasar por change

## Conventional Commits

Este proyecto sigue el estándar [Conventional Commits](https://www.conventionalcommits.org/) para mensajes de commit.

### Formato

```
<tipo>(<alcance>): <descripción>

[ cuerpo opcional ]

[ footer opcional ]
```

### Tipos

| Tipo | Descripción |
|------|-------------|
| `feat` | Nueva funcionalidad |
| `fix` | Corrección de bug |
| `docs` | Documentación |
| `style` | Formateo (sin cambio de lógica) |
| `refactor` | Refactoring (sin cambio de comportamiento) |
| `test` | Tests |
| `chore` | Tareas de mantenimiento |
| `opsx` | Cambios en specs OpenSpec |

### Ejemplos

```bash
feat(webhook): add endpoint /test-webhook for simulation
fix(pipeline): resolve nil pointer in normalization stage
docs(readme): update architecture diagram
opsx(specs): add tts-provider spec
```

## Worklogs

Registrar el trabajo realizado periódicamente en `logs/worklogs/`.

### Formato del Archivo

```
WORKLOG.<timestampunix>.md
```

Ejemplo: `WORKLOG.1741970000.md` (timestamp Unix)

### Estructura del Worklog

```markdown
---
title: "<título de la épica>"
summary: "<resumen corto de lo logrado>"
description: "<descripción más detallada>"
createdAt: "<fecha ISO 8601>"
tags:
  - tag1
  - tag2
metadata:
  epic: "<nombre-de-épica>"
  status: "completed|in-progress"
  repository: "whatsapp-tts"
  language: "es-ES"
---

# Épica: <nombre>

## Objetivo

[Qué se quería lograr]

## Componentes Implementados

- [ lista de componentes ]

## Problemas Encontrados y Soluciones

### Problema N: <título>
**Descripción:** [qué pasó]
**Solución:** [cómo se resolvió]

## Estado

- [ estado final ]

## Siguiente

[ qué sigue ]
```

### Cuándo Crear Worklog

- Al completar una épica o feature significativo
- Al resolver problemas importantes
- Al final de cada sesión de trabajo grande
- Al'archivar un change en OpenSpec
