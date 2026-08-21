type AuthAlertProps = {
  message: string;
};

export function AuthAlert({ message }: AuthAlertProps) {
  return (
    <div
      role="alert"
      className="rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-950"
    >
      {message}
    </div>
  );
}
