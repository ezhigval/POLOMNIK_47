/** Keep a supplied geographic area code (some Russian landlines have five digits). */
export function russianPhone(input: string) {
  const raw = input.trim();
  const [main, extension] = raw.split(/\s*(?:доб\.?|ext\.?)\s*/i, 2);
  const digits = main.replace(/\D/g, "");
  const national = digits.length === 11 && /^[78]/.test(digits) ? digits.slice(1)
    : digits.length === 10 && !main.startsWith("+") ? digits : null;

  if (!national || !/^[+\d()\s.-]+$/.test(main) || (extension && !/^\d+$/.test(extension))) {
    return { label: raw, href: "tel:" + raw.replace(/[()\s-]/g, "") };
  }
  const explicit = main.match(/\((\d{3,5})\)/)?.[1];
  const area = explicit && national.startsWith(explicit) ? explicit :
    national.match(/^(81367|81364)/)?.[1] ?? national.slice(0, 3);
  const subscriber = national.slice(area.length);
  const label = "+7 (" + area + ") " + subscriber.slice(0,-4) + "-" +
    subscriber.slice(-4,-2) + "-" + subscriber.slice(-2);
  return {
    label: label + (extension ? " доб. " + extension : ""),
    href: "tel:+7" + national + (extension ? ";ext=" + extension : ""),
  };
}
