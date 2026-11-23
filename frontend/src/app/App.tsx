import { CalculatorPage } from '../features/calculator/components/CalculatorPage'
import { AppErrorBoundary } from './AppErrorBoundary'

export function App() {
  return (
    <AppErrorBoundary>
      <CalculatorPage />
    </AppErrorBoundary>
  )
}
