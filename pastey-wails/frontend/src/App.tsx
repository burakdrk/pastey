import { ThemeProvider } from "./components/theme-provider";
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import { globalState } from "./lib/store";
import { useAtom } from "jotai";
import Home from "./pages/Home";
import Authentication from "./pages/Authentication";

function App() {
  const [isLoggedIn, setIsLoggedIn] = useAtom(globalState.isLoggedIn);

  return (
    <ThemeProvider defaultTheme="dark" storageKey="ui-theme">
      <Router>
        <Routes>
          <Route path="/" element={isLoggedIn ? <Home /> : <Navigate to="/login" />} />
          <Route path="/login" element={isLoggedIn ? <Navigate to="/" /> : <Authentication />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
}

export default App;
