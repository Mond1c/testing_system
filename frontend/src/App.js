import './App.css';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import FileUploader from "./components/FileUploader";
import Main from "./components/Main";

function App() {
  return (
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Main/>}/>
          <Route path="/upload" element={<FileUploader/>}/>
      </Routes>
      </BrowserRouter>
  );
}

export default App;
