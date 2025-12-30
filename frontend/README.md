# Kanban Board - Frontend

Modern Kanban Board built with Next.js 14, TypeScript, and Tailwind CSS.

## Tech Stack
- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript
- **State Management**: Zustand
- **UI Framework**: Tailwind CSS + shadcn/ui
- **HTTP Client**: Axios
- **Form Management**: React Hook Form + Zod
- **Architecture**: Atomic Design

## Project Structure

```
frontend/
├── app/                      # Next.js App Router
│   ├── (auth)/              # Auth group routes
│   │   ├── login/
│   │   └── register/
│   ├── (dashboard)/         # Dashboard group routes
│   │   ├── boards/
│   │   └── tasks/
│   ├── layout.tsx
│   ├── page.tsx
│   └── globals.css
├── src/
│   ├── components/          # Atomic Design Components
│   │   ├── atoms/          # Basic building blocks
│   │   ├── molecules/      # Simple component groups
│   │   ├── organisms/      # Complex components
│   │   └── templates/      # Page templates
│   ├── modules/            # Feature modules
│   │   ├── auth/
│   │   ├── board/
│   │   └── task/
│   ├── hooks/              # Custom React hooks
│   ├── services/           # API services
│   ├── store/              # Zustand stores
│   ├── utils/              # Utility functions
│   ├── types/              # TypeScript types
│   ├── styles/             # Global styles
│   └── config/             # Configuration
├── public/                 # Static assets
└── package.json
```

## Getting Started

1. **Install dependencies**
```bash
npm install
```

2. **Setup environment variables**
```bash
cp .env.example .env.local
```

3. **Run development server**
```bash
npm run dev
```

4. **Build for production**
```bash
npm run build
npm start
```

## Available Scripts
- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint
- `npm run type-check` - Run TypeScript type checking

## Atomic Design Structure

### Atoms
Basic UI elements: Button, Input, Label, Badge, etc.

### Molecules
Simple component combinations: FormField, SearchBar, Card, etc.

### Organisms
Complex components: Navbar, TaskCard, BoardList, etc.

### Templates
Page layouts: DashboardTemplate, AuthTemplate, etc.

## API Integration

All API calls are centralized in `src/services/`:
- `authService.ts` - Authentication endpoints
- `userService.ts` - User management
- `boardService.ts` - Board operations
- `taskService.ts` - Task operations

## State Management

Using Zustand for global state:
- `authStore.ts` - Authentication state
- `boardStore.ts` - Board state
- `taskStore.ts` - Task state
- `uiStore.ts` - UI state (modals, toasts, etc.)

## Environment Variables

See `.env.example` for required variables.
