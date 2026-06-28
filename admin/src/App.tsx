import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './store/auth'
import MainLayout from './components/layout/MainLayout'
import LoginPage from './pages/auth/LoginPage'
import DashboardPage from './pages/dashboard/DashboardPage'
import UserListPage from './pages/users/UserListPage'
import UserDetailPage from './pages/users/UserDetailPage'
import AdminListPage from './pages/admins/AdminListPage'
import AdminFormPage from './pages/admins/AdminFormPage'
import RoleListPage from './pages/roles/RoleListPage'
import RoleFormPage from './pages/roles/RoleFormPage'
import ReportsPage from './pages/reports/ReportsPage'
import ExportPage from './pages/reports/ExportPage'
import LogsPage from './pages/logs/LogsPage'
import MediaContentPage from './pages/media/MediaContentPage'
import ProfilePage from './pages/profile/ProfilePage'
import CalendarsPage from './pages/data/CalendarsPage'
import EmotionsPage from './pages/data/EmotionsPage'
import StressEventsPage from './pages/data/StressEventsPage'
import BodyTensionPage from './pages/data/BodyTensionPage'
import BreathingPage from './pages/data/BreathingPage'
import CognitiveGamesPage from './pages/data/CognitiveGamesPage'
import MentalMustsPage from './pages/data/MentalMustsPage'
import NegativeThoughtsPage from './pages/data/NegativeThoughtsPage'
import MindCourtPage from './pages/data/MindCourtPage'
import ConflictExercisesPage from './pages/data/ConflictExercisesPage'
import MoodTrackerPage from './pages/data/MoodTrackerPage'
import RolesValuesPage from './pages/data/RolesValuesPage'
import SkyThoughtsPage from './pages/data/SkyThoughtsPage'
import MindfulTimersPage from './pages/data/MindfulTimersPage'
import AcceptancePage from './pages/data/AcceptancePage'
import WeeklyReportsPage from './pages/data/WeeklyReportsPage'

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const { isAuthenticated } = useAuthStore()

    if (!isAuthenticated) {
        return <Navigate to="/login" replace />
    }

    return <MainLayout>{children}</MainLayout>
}

function App() {
    return (
        <Routes>
            { }
            <Route path="/login" element={<LoginPage />} />

            { }
            <Route
                path="/"
                element={
                    <ProtectedRoute>
                        <DashboardPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/users"
                element={
                    <ProtectedRoute>
                        <UserListPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/users/:id"
                element={
                    <ProtectedRoute>
                        <UserDetailPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/admins"
                element={
                    <ProtectedRoute>
                        <AdminListPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/admins/new"
                element={
                    <ProtectedRoute>
                        <AdminFormPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/admins/:id/edit"
                element={
                    <ProtectedRoute>
                        <AdminFormPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/roles"
                element={
                    <ProtectedRoute>
                        <RoleListPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/roles/new"
                element={
                    <ProtectedRoute>
                        <RoleFormPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/roles/:id/edit"
                element={
                    <ProtectedRoute>
                        <RoleFormPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/reports"
                element={
                    <ProtectedRoute>
                        <ReportsPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/reports/export"
                element={
                    <ProtectedRoute>
                        <ExportPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/logs"
                element={
                    <ProtectedRoute>
                        <LogsPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/media"
                element={
                    <ProtectedRoute>
                        <MediaContentPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/profile"
                element={
                    <ProtectedRoute>
                        <ProfilePage />
                    </ProtectedRoute>
                }
            />

            { }
            <Route
                path="/data/calendars"
                element={
                    <ProtectedRoute>
                        <CalendarsPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/data/emotions"
                element={
                    <ProtectedRoute>
                        <EmotionsPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/data/stress-events"
                element={
                    <ProtectedRoute>
                        <StressEventsPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/data/body-tension"
                element={
                    <ProtectedRoute>
                        <BodyTensionPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/data/breathing"
                element={
                    <ProtectedRoute>
                        <BreathingPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/data/cognitive-games"
                element={
                    <ProtectedRoute>
                        <CognitiveGamesPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/data/mental-musts"
                element={
                    <ProtectedRoute>
                        <MentalMustsPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/data/negative-thoughts"
                element={
                    <ProtectedRoute>
                        <NegativeThoughtsPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/data/mind-court"
                element={
                    <ProtectedRoute>
                        <MindCourtPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/data/conflict-exercises"
                element={
                    <ProtectedRoute>
                        <ConflictExercisesPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/data/mood-tracker"
                element={
                    <ProtectedRoute>
                        <MoodTrackerPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/data/roles-values"
                element={
                    <ProtectedRoute>
                        <RolesValuesPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/data/sky-thoughts"
                element={
                    <ProtectedRoute>
                        <SkyThoughtsPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/data/mindful-timers"
                element={
                    <ProtectedRoute>
                        <MindfulTimersPage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/data/acceptance"
                element={
                    <ProtectedRoute>
                        <AcceptancePage />
                    </ProtectedRoute>
                }
            />
            <Route
                path="/data/weekly-reports"
                element={
                    <ProtectedRoute>
                        <WeeklyReportsPage />
                    </ProtectedRoute>
                }
            />

            { }
            <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
    )
}

export default App
