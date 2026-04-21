# Skill: Next.js & Frontend Core

This skill defines the foundational standards for developing the frontend within the AI Agency Platform. It focuses on architecture, design constraints, and proper interaction with the backend.

## Standards

### 1. Architecture & Tech Stack
- **Framework**: Use **Next.js (App Router)** with **React**, **TypeScript**, and **TailwindCSS**.
- **UI Components**: Use **shadcn** (ensure the package is `shadcn`, not the deprecated `shadcn-ui`), `lucide-react`, `react-hook-form`, and `zod`.
- **Drag and Drop**: Use `@dnd-kit` for drag-and-drop interfaces.
- **Data Fetching**: The application integrates directly with a fully functional backend API via Axios. **Never use mock data.**

### 2. Design & Visual Identity
- **Theme**: Strict **Dark Mode** theme.
- **Color Palette**: Must consist of **Dark Gray**, **Pool Blue**, and **Water Green**.
- **UX/UI**: Prioritize excellent layout, UX, and usability, as this frontend serves as a showcase for a software development agency.

### 3. Authentication & State
- **State Management**: Use React Context API for state management.
- **Tokens**:
  - Access tokens MUST be stored in `localStorage`.
  - Refresh tokens MUST be handled via HttpOnly cookies managed by the backend.
- **Route Protection**: Use client-side route protection (e.g., via a `PrivateRoute` component). Do not use Next.js middleware for this due to `localStorage` limitations.
- **Interceptor Detail**: In the `api.ts` Axios interceptor, use `axios.post` directly (instead of the custom `api` instance) for refresh token requests to avoid infinite loops and deadlocks.

### 4. Component Implementation Rules
- **Shadcn `<Button>`**: Avoid using the `asChild` prop on shadcn `<Button>` components. It causes Next.js TypeScript compiler type-checking errors (`IntrinsicAttributes & ButtonProps`) with the current configuration.

## Examples

### Good Authentication Setup
```tsx
// Using PrivateRoute for client-side protection instead of Next.js Middleware
export default function DashboardPage() {
  return (
    <PrivateRoute>
      <DashboardContent />
    </PrivateRoute>
  );
}
```

### Good Axios Interceptor Setup
```typescript
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    // ... setup to detect token expiration ...
    // GOOD: using axios.post directly to prevent deadlocks
    const response = await axios.post('/api/refresh-token', {}, { withCredentials: true });
    // ... apply new token ...
  }
);
```
