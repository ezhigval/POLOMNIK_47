type ManagementPanelProps = {
  title?: string;
  description?: string;
  children: React.ReactNode;
  className?: string;
};

export function ManagementPanel({ title, description, children, className = "" }: ManagementPanelProps) {
  return (
    <section className={`rounded-2xl border border-stone-200 bg-white ${className}`}>
      {title ? (
        <div className="border-b border-stone-100 px-5 py-4">
          <h2 className="text-base font-semibold text-stone-900">{title}</h2>
          {description ? <p className="mt-1 text-sm text-stone-600">{description}</p> : null}
        </div>
      ) : null}
      {children}
    </section>
  );
}

export function ManagementTable({ children }: { children: React.ReactNode }) {
  return (
    <div className="overflow-x-auto">
      <table className="min-w-full text-left text-sm">{children}</table>
    </div>
  );
}

export function ManagementTableHead({ children }: { children: React.ReactNode }) {
  return (
    <thead className="border-b border-stone-200 bg-stone-50 text-stone-600">
      <tr>{children}</tr>
    </thead>
  );
}

export function ManagementTh({ children }: { children?: React.ReactNode }) {
  return <th className="px-4 py-3 font-medium">{children}</th>;
}

export function ManagementEmptyRow({ colSpan, children }: { colSpan: number; children: React.ReactNode }) {
  return (
    <tr>
      <td colSpan={colSpan} className="px-4 py-12 text-center text-stone-500">
        {children}
      </td>
    </tr>
  );
}
