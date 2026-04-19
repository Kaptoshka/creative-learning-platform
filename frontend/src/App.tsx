import { Routes, Route } from "react-router-dom";
import HomePage from "@/pages/HomePage";
import DashboardPage from "@/pages/DashboardPage";
import TasksPage from "@/pages/TasksPage";
import TaskPage from "@/pages/TaskPage";
import CreateAssignmentPage from "@/pages/CreateAssignmentPage";
import ReviewPage from "@/pages/ReviewPage";
import SubmissionsPage from "@/pages/SubmissionsPage";
import AuthPage from "@/pages/AuthPage";
import TeacherDashboardPage from "@/pages/TeacherDashboardPage";
import TaskDetailPage from "@/pages/TaskDetailPage";
import { AuthProvider } from "@/context/AuthProvider";
import Navigation from "@/components/Navigation";
import ProtectedRoutes from "@/components/ProtectedRoutes";
import "@/index.css";

function App() {
    return (
        <AuthProvider>
            <Navigation />
            <Routes>
                <Route path="/" element={<HomePage />} />
                <Route path="/auth" element={<AuthPage />} />
                <Route path="/tasks" element={<TasksPage />} />
                <Route path="/task/:id" element={<TaskPage />} />
                <Route path="/task-detail/:id" element={<TaskDetailPage />} />
                <Route element={<ProtectedRoutes />}>
                    <Route path="/dashboard" element={<DashboardPage />} />
                    <Route
                        path="/create-assignment"
                        element={<CreateAssignmentPage />}
                    />
                    <Route path="/review" element={<ReviewPage />} />
                    <Route path="/submissions" element={<SubmissionsPage />} />
                    <Route
                        path="/teacher-dashboard"
                        element={<TeacherDashboardPage />}
                    />
                </Route>
            </Routes>
        </AuthProvider>
    );
}

export default App;
