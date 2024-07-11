import { ThemeProvider } from "./components/theme-provider";
import { HashRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import { globalState } from "./lib/store";
import { useAtom } from "jotai";
import Root from "./pages/Root";
import Authentication from "./pages/Authentication";
import Header from "./components/header";
import { useEffect, useState } from "react";
import { GetIsLoggedIn } from "../wailsjs/go/backend/App";
import Splash from "./pages/Splash";
import { sleep } from "./lib/utils";

function App() {
  const [isLoggedIn, setIsLoggedIn] = useAtom(globalState.isLoggedIn);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function checkLogin() {
      const loggedIn = await GetIsLoggedIn();
      await sleep(1000);
      setIsLoggedIn(loggedIn);
      setLoading(false);
    }

    checkLogin();
  }, []);

  return (
    <ThemeProvider defaultTheme="dark" storageKey="ui-theme">
      <Header />
      {loading ? (
        <Splash />
      ) : (
        <Router>
          <Routes>
            <Route path="/" element={isLoggedIn ? <Root /> : <Navigate to="/login" />} />
            <Route path="/login" element={isLoggedIn ? <Navigate to="/" /> : <Authentication />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </Router>
      )}
    </ThemeProvider>
  );
}

export default App;
