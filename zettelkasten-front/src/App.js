import { Admin } from "./pages/AdminPage";
import LandingPage from "./pages/LandingPage";
import LoginForm from "./pages/LoginPage";
import MainApp from "./pages/MainApp";
import RegisterPage from "./pages/RegisterPage";
import { Routes, Route } from "react-router-dom";
import { ProtectedAdminPage } from "./components/ProtectedAdminPage";
import { AdminUserDetailPage } from "./pages/AdminUserDetailPage";
import { AdminEditUserPage } from "./pages/AdminEditUserPage";
import PasswordReset from "./pages/PasswordReset";

function App() {
  return (
    <div>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/app" element={<MainApp />} />
        <Route
          path="/admin"
          element={
            <ProtectedAdminPage>
              <Admin />
            </ProtectedAdminPage>
          }
        />
        <Route
          path="/admin/user/:id"
          element={
            <ProtectedAdminPage>
              <AdminUserDetailPage />
            </ProtectedAdminPage>
          }
        />
	  <Route
      path="/admin/user/:id/edit"
          element={
            <ProtectedAdminPage>
              <AdminEditUserPage />
            </ProtectedAdminPage>
          }
      />
        <Route path="/login" element={<LoginForm />} />
        <Route path="/register" element={<RegisterPage />} />
	  <Route path="/reset" element={<PasswordReset />} />
      </Routes>
    </div>
  );
}

export default App;
