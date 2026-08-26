type TourRegularCheckboxProps = {
  checked: boolean;
  onChange: (checked: boolean) => void;
};

export function TourRegularCheckbox({ checked, onChange }: TourRegularCheckboxProps) {
  return (
    <label className="flex items-start gap-2 text-sm md:col-span-2">
      <input
        type="checkbox"
        name="is_regular"
        checked={checked}
        onChange={(event) => onChange(event.target.checked)}
        className="mt-1 size-4"
      />
      <span>
        <span className="block font-medium">Регулярный тур</span>
        <span className="block text-stone-500">
          Скрывает цены и даты на сайте и помечает как регулярные туры.
        </span>
      </span>
    </label>
  );
}
