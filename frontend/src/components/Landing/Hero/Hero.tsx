import { useEffect, useState } from 'react'

import { useMemoizedUser } from '../../../hooks/useMemoizedUser'
import Button from '../../Button'

import {
  ContentContainer,
  Image,
  ImageContainer,
  Section,
  Title,
} from './Hero.styles'

import HeroImg from '/src/assets/images/Programming-rafiki.svg'

interface TypingAnimationProps {
  text: string
  typingSpeed: number
}

const TypingAnimation: React.FC<TypingAnimationProps> = ({
  text,
  typingSpeed,
}) => {
  const [displayedText, setDisplayedText] = useState('')
  const [isFinished, setIsFinished] = useState(false)

  useEffect(() => {
    let index = 0
    let accumulatedTime = 0
    const startTime = performance.now()
    let requestId: number

    const typeCharacter = (timestamp: number) => {
      if (index < text.length) {
        accumulatedTime += timestamp - startTime

        if (accumulatedTime >= typingSpeed) {
          setDisplayedText((prev) => prev + text.charAt(index))
          index++
          accumulatedTime = 0
        }

        requestId = requestAnimationFrame(typeCharacter)
      }

      if (index === text.length) {
        setTimeout(() => {
          setIsFinished(true)
        }, 1000)
      }
    }

    requestId = requestAnimationFrame(typeCharacter)

    return () => cancelAnimationFrame(requestId)
  }, [text, typingSpeed])

  return (
    <span className="typing-animation">
      {displayedText}
      {!isFinished && (
        <span className="ml-3 animate-ping text-gray-800">|</span>
      )}
    </span>
  )
}

const AchievementsSection = () => {
  const { user } = useMemoizedUser()
  return (
    <Section>
      <ContentContainer>
        <p className="mb-4 text-center font-semibold text-mainColor md:mb-6 md:text-lg lg:text-left xl:text-xl">
          Very proud to introduce
        </p>

        <Title>
          Achievements on <br />
          <TypingAnimation text="Thhe Open Network" typingSpeed={3000} />
          {/* <span>The Open Network</span> */}
        </Title>

        <p className="mb-8 text-center font-[500] leading-7 text-gray-700 lg:text-left">
          Tap into your secret powers with the TON Developers Platform, an
          exciting new world that turns your open-source contributions into SBT
          rewards, hooking you up with exclusive, real-world TON collectibles!
        </p>

        <div className="relative flex justify-center lg:justify-start">
          <div className="flex justify-center gap-2.5 lg:justify-start">
            {!user && (
              <>
                <Button href="https://docs.ton.org" color="blue">
                  Learn more
                </Button>{' '}
                <Button href="https://t.me/goreactdev" color="white">
                  Try Bot
                </Button>
              </>
            )}

            {user && (
              <>
                {' '}
                <Button to="/rewards" color="blue">
                  Explore members
                </Button>{' '}
              </>
            )}
          </div>
        </div>
      </ContentContainer>

      <ImageContainer>
        <Image src={HeroImg} loading="lazy" alt="Photo by" />
      </ImageContainer>
    </Section>
  )
}

export default AchievementsSection
