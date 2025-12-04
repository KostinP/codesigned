'use client'
import { Label } from '@/components/ui/label'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'

const roles = [
  'Разработчик',
  'Дизайнер',
  'Геймдизайнер',
  'Менеджер продукта',
  'Тестировщик',
  'Контент-мейкер',
  'Аналитик',
]

export default function StepRole({ value, onChange }: { value: string, onChange: (val: string) => void }) {
  return (
    <div>
      <Label className="mb-4 block">Какая роль больше всего описывает вас?</Label>
      <RadioGroup value={value} onValueChange={onChange} className="space-y-2">
        {roles.map(r => (
          <div key={r} className="flex items-center space-x-2">
            <RadioGroupItem value={r} id={r} />
            <Label htmlFor={r}>{r}</Label>
          </div>
        ))}
      </RadioGroup>
    </div>
  )
}
