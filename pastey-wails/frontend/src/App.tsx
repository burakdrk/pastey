import { ThemeProvider } from "./components/theme-provider";
import { HashRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import { globalState } from "./lib/store";
import { useAtom } from "jotai";
import Root from "./pages/Root";
import Authentication from "./pages/Authentication";
import Header from "./components/header";
import { useEffect } from "react";
import { GetIsLoggedIn } from "../wailsjs/go/backend/App";

function App() {
  const [isLoggedIn, setIsLoggedIn] = useAtom(globalState.isLoggedIn);

  useEffect(() => {
    async function checkLogin() {
      const loggedIn = await GetIsLoggedIn();
      setIsLoggedIn(loggedIn);
    }

    checkLogin();
  }, []);

  return (
    <ThemeProvider defaultTheme="dark" storageKey="ui-theme">
      <Header />
      <Router>
        <Routes>
          <Route path="/" element={isLoggedIn ? <Root /> : <Navigate to="/login" />} />
          <Route path="/login" element={isLoggedIn ? <Navigate to="/" /> : <Authentication />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
}

export default App;
