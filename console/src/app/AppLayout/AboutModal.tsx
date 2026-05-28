import { AboutModal, Content } from '@patternfly/react-core';
import React from 'react';
import brandImg from '../../assets/favicon.png';
import modalBackground from '../../assets/modal_background.png';
import { APP_VERSION, PRODUCT_NAME, MAINTAINER_NAME, REPOSITORY_URL } from '@app/constants';

interface AboutModalComponentProps {
  isOpen: boolean;
  onClose: () => void;
}

const AboutModalComponent: React.FunctionComponent<AboutModalComponentProps> = ({ isOpen, onClose }) => {
  return (
    <AboutModal
      isOpen={isOpen}
      onClose={onClose}
      productName={PRODUCT_NAME}
      brandImageAlt="Red Hat logo"
      brandImageSrc={brandImg}
      backgroundImageSrc={modalBackground}
    >
      <Content>
        <Content component="dl">
          <Content component="dt">Version</Content>
          <Content component="dd">{APP_VERSION}</Content>

          <Content component="dt">Maintainer</Content>
          <Content component="dd">{MAINTAINER_NAME}</Content>

          <Content component="dt">Repository</Content>
          <Content component="dd">
            <a href={REPOSITORY_URL} target="_blank" rel="noopener noreferrer">
              GitHub
            </a>
          </Content>
        </Content>
      </Content>
    </AboutModal>
  );
};

export default AboutModalComponent;
