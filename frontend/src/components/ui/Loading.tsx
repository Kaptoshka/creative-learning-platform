interface LoadingProps {
  text?: string;
  fullPage?: boolean;
}

const Loading = ({ text = "Загрузка...", fullPage = false }: LoadingProps) => {
  const containerClass = fullPage ? "page-container" : "";

  return (
    <div className={containerClass}>
      <div className="loading">
        <div className="loading__spinner"></div>
        <span className="loading__text">{text}</span>
      </div>
    </div>
  );
};

export default Loading;
