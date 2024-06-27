import "./App.css";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import FileUploader from "./components/FileUploader";
import Main from "./components/Main";
import Results from "./components/Results";
import Runs from "./components/Runs";
import AllRuns from "./components/AllRuns";
import Header from "./components/Header";
import Problems from "./components/Problems";

function App() {
  return (
    <BrowserRouter>
      <div>
        <Header/>
        <Routes>
          <Route path="/" element={<Main />} />
          <Route path="/upload" element={<FileUploader />} />
          <Route path="/results" element={<Results />} />
          <Route path="/runs" element={<Runs />} />
          <Route path="/admin/all_runs" element={<AllRuns />} />
          <Route path="/problems" element={<Problems />} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;
