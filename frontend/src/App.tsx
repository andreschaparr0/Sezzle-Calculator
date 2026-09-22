import { Calculator } from './components/Calculator';

export default function App() {
  return (
    <main className="app">
      <header className="app__header">
        <h1>Calculator</h1>
        <p>Arithmetic is performed by the backend service.</p>
      </header>
      <Calculator />
    </main>
  );
}
