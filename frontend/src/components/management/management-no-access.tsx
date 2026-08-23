export function ManagementNoAccess() {
  return (
    <p className="rounded-2xl border border-stone-200 bg-white p-5 text-sm text-stone-600">
      У этой роли нет доступа к разделу. В меню остаются только пункты по правам; API по-прежнему
      проверяет те же разрешения.
    </p>
  );
}
