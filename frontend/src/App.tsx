import { Calculator } from './components/Calculator';

export default function App() {
  return (
    <main className="app">
      <div className="app__badge">
        <img
          className="app__logo"
          src="/sezzle-logo.jpg"
          alt="Sezzle"
          onError={(event) => {
            event.currentTarget.style.display = 'none';
          }}
        />
        <div className="app__badge-text">
          <span className="app__badge-title">Technical assessment</span>
          <span className="app__badge-name">Andres Felipe Chaparro Diaz</span>
        </div>
      </div>

      <header className="app__header">
        <h1>Calculator</h1>
        <p>Arithmetic is performed by the backend service.</p>
      </header>
      <Calculator />
    </main>
  );
}
