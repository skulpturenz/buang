import { Route, Routes } from "react-router-dom";
import { ProtectedRoute } from "./components/ProtectedRoute.js";
import { CompleteProfile } from "./pages/CompleteProfile.js";
import { Login } from "./pages/Login.js";
import { Register } from "./pages/Register.js";
import { Setup } from "./pages/Setup.js";

export const App = () => (
	<Routes>
		<Route path="/login" element={<Login />} />
		<Route path="/register" element={<Register />} />
		<Route path="/complete-profile" element={<CompleteProfile />} />
		<Route
			path="/setup"
			element={
				<ProtectedRoute>
					<Setup />
				</ProtectedRoute>
			}
		/>
		<Route path="*" element={<Login />} />
	</Routes>
);
