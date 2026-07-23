import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { Layout } from './components/layout/Layout'
import { HomePage } from './pages/HomePage'
import { CompanyDashboardPage } from './features/company-dashboard/pages/CompanyDashboardPage'
import { NewsChartsPage } from './features/news-charts/pages/NewsChartsPage'
import { LearningCenterPage } from './features/learning-center/pages/LearningCenterPage'
import { PortfolioPage } from './features/portfolio/pages/PortfolioPage'
import { LoginPage } from './features/auth/pages/LoginPage'
import { SignUpPage } from './features/auth/pages/SignUpPage'
import { RequireGuest } from './features/auth/components/RequireGuest'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/login"
          element={
            <RequireGuest>
              <LoginPage />
            </RequireGuest>
          }
        />
        <Route
          path="/register"
          element={
            <RequireGuest>
              <SignUpPage />
            </RequireGuest>
          }
        />
        <Route
          path="/*"
          element={
            <Layout>
              <Routes>
                <Route path="/" element={<HomePage />} />
                <Route path="/dashboard" element={<CompanyDashboardPage />} />
                <Route path="/news" element={<NewsChartsPage />} />
                <Route path="/learn" element={<LearningCenterPage />} />
                <Route path="/portfolio" element={<PortfolioPage />} />
              </Routes>
            </Layout>
          }
        />
      </Routes>
    </BrowserRouter>
  )
}

export default App
