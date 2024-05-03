import "./App.css";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import FileUploader from "./components/FileUploader";
import Main from "./components/Main";
import Results from "./components/Results";
import Runs from "./components/Runs";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Main />} />
        <Route path="/upload" element={<FileUploader />} />
        <Route path="/results" element={<Results />} />
        <Route path="/runs" element={<Runs/>} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
