import { arrayOfIcons } from '../Icons/Icons'
import { SocialLink, SocialSvg } from '../Icons/Icons.styles'

import {
  GridContainer,
  LogoContainer,
  LogoLink,
  LogoSvg,
  NavHeading,
  NavLink,
  NavWrapper,
  Text,
} from './Footer.styles'
import { footerLinks } from './links'

const Footer = () => {
  return (
    <div>
      <GridContainer>
        <div className="col-span-full lg:col-span-2">
          <LogoContainer>
            <LogoLink href="/" aria-label="logo">
              <LogoSvg
                width="95"
                height="94"
                viewBox="0 0 95 94"
                fill="currentColor"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path d="M96 0V47L48 94H0V47L48 0H96Z" />
              </LogoSvg>
              TDP
            </LogoLink>
          </LogoContainer>

          <Text>
            Filler text is dummy text which has no meaning however looks very
            similar to real text.
          </Text>

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
        </div>

        {footerLinks.map((footerLink) => (
          <div key={footerLink.heading}>
            <NavHeading>{footerLink.heading}</NavHeading>
            <NavWrapper>
              {footerLink.links.map((link) => (
                <div key={link.name}>
                  <NavLink href={link.url}>{link.name}</NavLink>
                </div>
              ))}
            </NavWrapper>
          </div>
        ))}
      </GridContainer>

      <div className="border-t py-8 text-center text-sm text-gray-400">
        © 2023 - Present TON Developers Platform. All rights reserved.
      </div>
    </div>
  )
}

export default Footer
