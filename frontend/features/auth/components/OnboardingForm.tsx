'use client'
import { useState } from 'react'
import StepRole from './StepRole'
import StepEmail from './StepEmail'
import { Button } from '@/components/ui/button'

export default function OnboardingForm() {
  const [step, setStep] = useState(0)
  const [role, setRole] = useState('')
  const [email, setEmail] = useState('')
  const [subscribe, setSubscribe] = useState(false)

  const next = () => setStep(s => s + 1)
  const prev = () => setStep(s => Math.max(0, s - 1))

  const handleSubmit = async () => {
    // await updateUser({
    //   name: 'Имя по умолчанию', // временно
    //   role,
    //   email: subscribe ? email : undefined,
    //   subscribe,
    // })
    // window.location.href = '/'
  }

  return (
    <div className="space-y-6">
      {step === 0 && <StepRole value={role} onChange={setRole} />}
      {step === 1 && <StepEmail
        email={email}
        subscribed={subscribe}
        onEmailChange={setEmail}
        onSubscribeChange={setSubscribe}
      />}

      <div className="flex justify-between">
        {step > 0 && <Button variant="outline" onClick={prev}>Назад</Button>}
        {step < 1 ? (
          <Button onClick={next} disabled={!role}>Далее</Button>
        ) : (
          <Button onClick={handleSubmit} disabled={!role}>Завершить</Button>
        )}
      </div>
    </div>
  )
}
