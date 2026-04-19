interface ErrorMessageProps {
  error: string;
  onRetry?: () => void;
}

const ErrorMessage = ({ error, onRetry }: ErrorMessageProps) => {
  return (
    <div className="error-message">
      <p>{error}</p>
      {onRetry && (
        <button onClick={onRetry} className="error-message__retry">
          Повторить
        </button>
      )}
    </div>
  );
};

export default ErrorMessage;
