import { useEffect, useState } from 'react'

import { arrayOfIcons } from '../../Icons/Icons'
import { SocialLink, SocialSvg } from '../../Icons/Icons.styles'

import {
  Divider,
  Section,
  SocialContainer,
  SocialLabel,
  Stat,
  StatLabel,
  StatNumber,
  StatsContainer,
} from './Stats.styles'

interface AnimatedCounterProps {
  start: number
  end: number
  duration: number
}

const AnimatedCounter: React.FC<AnimatedCounterProps> = ({
  start,
  end,
  duration,
}) => {
  const [count, setCount] = useState(start)

  useEffect(() => {
    const startTime = Date.now()

    const animate = () => {
      const now = Date.now()
      const elapsed = now - startTime
      const progress = Math.min(elapsed / duration, 1)
      const currentCount = start + progress * (end - start)

      setCount(Math.round(currentCount))

      if (progress < 1) {
        requestAnimationFrame(animate)
      }
    }

    requestAnimationFrame(animate)
  }, [start, end, duration])

  return <div className="animate-[pulse_1s_ease-in-out]">{count}+</div>
}

const StatsSection = () => {
  return (
    <Section>
      <StatsContainer>
        <Stat>
          <StatNumber>
            <AnimatedCounter start={1} end={300} duration={500} />
          </StatNumber>
          <StatLabel>Developers</StatLabel>
        </Stat>

        <Stat>
          <StatNumber>
            <AnimatedCounter start={1} end={500} duration={500} />
          </StatNumber>
          <StatLabel>Projects</StatLabel>
        </Stat>

        <Stat>
          <StatNumber>
            <AnimatedCounter start={1} end={250} duration={500} />
          </StatNumber>
          <StatLabel>Rewards</StatLabel>
        </Stat>

        <Stat>
          <StatNumber>A couple</StatNumber>
          <StatLabel>Coffee breaks</StatLabel>
        </Stat>
      </StatsContainer>

      <SocialContainer>
        <SocialLabel>Social</SocialLabel>
        <Divider />

        <div className="flex gap-4">
          {arrayOfIcons.map((icon) => (
            <SocialLink key={icon.name} href="#" target="_blank">
              <SocialSvg
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="currentColor"
                xmlns="http://www.w3.org/2000/svg"
              >
                {icon.jsx}
              </SocialSvg>
            </SocialLink>
          ))}
        </div>
      </SocialContainer>
    </Section>
  )
}

export default StatsSection
