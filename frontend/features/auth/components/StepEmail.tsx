'use client'
import { Input } from '@/components/ui/input'
import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'

export default function StepEmail({
  email,
  subscribed,
  onEmailChange,
  onSubscribeChange,
}: {
  email: string,
  subscribed: boolean,
  onEmailChange: (val: string) => void,
  onSubscribeChange: (val: boolean) => void,
}) {
  return (
    <div className="space-y-4">
      <div className="flex items-center space-x-2">
        <Checkbox id="subscribe" checked={subscribed} onCheckedChange={onSubscribeChange} />
        <Label htmlFor="subscribe">Хочу получать рассылку по почте</Label>
      </div>

      {subscribed && (
        <div>
          <Label htmlFor="email">Почта</Label>
          <Input
            id="email"
            type="email"
            placeholder="example@email.com"
            value={email}
            onChange={e => onEmailChange(e.target.value)}
          />
        </div>
      )}
    </div>
  )
}
