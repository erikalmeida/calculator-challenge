import { Component, type ErrorInfo, type ReactNode } from 'react'

interface AppErrorBoundaryProps {
  children: ReactNode
}

interface AppErrorBoundaryState {
  hasError: boolean
}

export class AppErrorBoundary extends Component<AppErrorBoundaryProps, AppErrorBoundaryState> {
  state: AppErrorBoundaryState = { hasError: false }

  static getDerivedStateFromError(): AppErrorBoundaryState {
    return { hasError: true }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    if (import.meta.env.DEV) {
      console.error('Unexpected application error', { error, componentStack: errorInfo.componentStack })
    }
  }

  render() {
    if (this.state.hasError) {
      return (
        <main className="calculator-page">
          <section className="calculator-card" role="alert">
            <p className="calculator-eyebrow">Sezzle take-home</p>
            <h1>Something went wrong</h1>
            <p className="calculator-description">Please refresh the page and try again.</p>
          </section>
        </main>
      )
    }

    return this.props.children
  }
}
