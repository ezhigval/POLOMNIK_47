type FormErrorProps = {
  children?: string | null;
  className?: string;
};

export function FormError({ children, className = "" }: FormErrorProps) {
  if (!children) {
    return null;
  }

  return <p className={`text-sm text-red-700 ${className}`.trim()}>{children}</p>;
}
