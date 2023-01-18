import { createDialog } from "../../components/dialog";
import styled from "styled-components";
import React, { useContext, useState } from "react";
import { Qrcode } from "../../components/qrcode";
import { Loading } from "../../components/loading";
import { AppContext } from "../../shared/app_context";
import { Center } from "../../components/center";
import { http } from "../../shared/http";
import { getWsClient } from "../../shared/ws_client";
import { clientId } from "../../initializers/client_id";
import { Space } from "../../components/space";

export const Layout = styled.div`
  min-height: 100vh; display: flex; align-items: stretch; flex-direction: column;
  padding: 0 16px; margin: 0 auto;
  @media (min-width: 414px) {
    max-width: 600px; 
  }
`;
export const Header = styled.h1`
  margin-top: 48px;
  margin-bottom: 32px;
  text-align: center;
`;
export const BigTextarea = styled.textarea`
  width: 100%;
  min-height: 160px;
  line-height: 20px;
  &.draging {
    border-color: red;
  }
`;
export const Button = styled.button`
  min-height: 40px;
  padding: 0 60px;
  border: 2px solid #666;
  background: #f5b70d;
  border-radius: 8px;
  cursor: pointer;
`;
export const Form = styled.form`
  &> .row {
    margin: 16px 0;
  }
`;

const Span = styled.span`
  margin-right: 8px;
`;
const Label = styled.label`
  display: flex; padding: 4px 0; 
  justify-content: flex-start; align-items: center;
  min-height: 40px; white-space: nowrap;
`;
const H2 = styled.h2`
  font-weight: bold; font-size: 24px;
  margin-bottom: 16px;
`
const P = styled.p`
  a {text-decoration: underline;}
`
const UploadSuccessDialog = ({ content, onClose }) => {
  const [address, setAddress] = useState(localStorage.getItem("address") || "");
  const context = useContext(AppContext);
  const addressesRef = context?.addressesRef ?? null
  const onChange = (e) => {
    setAddress(e.target.value);
    localStorage.setItem("address", e.target.value);
  };
  content = typeof content === "string" ? content : content(address)
  return (
    <Pop>
      <H2>上傳成功</H2>
      {addressesRef.current ?
        <div>
          <P>
            請 Windows 使用者在防火牆入站規則開啟 27149 端口（<a href="https://jingyan.baidu.com/article/09ea3ede7311dec0afde3977.html" target="_blank" rel="noreferrer">教程</a>）
          </P>
          <P>
            <Label>
              <Span>請選擇手機可以訪問的區域網IP</Span>
              <select value={address} onChange={onChange}>
                <option value="" disabled>
                  - 請選擇 -
                </option>
                {addressesRef.current?.map((string) => (
                  <option key={string}>{string}</option>
                ))}
              </select>
            </Label>
          </P>
        </div>
        : null
      }
      <Center>
        {content ? <Qrcode content={content} /> : null}
      </Center>
      <Center>
        {content ? <a href={content}>請 手機掃碼 或 點擊下載</a> : null}
      </Center>
      <Space />
      <Center>
        <Button onClick={onClose}>關閉</Button>
      </Center>
    </Pop>
  );
};
export const showUploadTextSuccessDialog = ({ context, content }) => {
  const close = createDialog(
    <UploadSuccessDialog content={content} onClose={() => close()} />,
    { context }
  );
};
export const showUploadFileSuccessDialog = ({ context, content }) => {
  const close = createDialog(
    <UploadSuccessDialog content={content} onClose={() => close()} />,
    { context }
  );
};
export const showUploadFailDialog = () => {
  return createDialog(
    <Pop>
      <div>上傳失敗</div>
      <button onClick={() => close()}>關閉</button>
    </Pop>
  );
};
export const showUploadingDialog = () => {
  return createDialog(<Loading>上傳中</Loading>);
};
const Pop = styled.div`
  padding: 16px;
`;

const notifyPc = (response, type) => {
  getWsClient().then(c => {
    c.send({ clientId, type, url: response.data.url })
  })
  return response
}
export const uploadText = (text) => {
  return http.post("/api/v1/texts", {
    raw: text
  }).then(r => notifyPc(r, 'text'))
}
export const uploadFile = (blob) => {
  const formData = new FormData();
  formData.append("raw", blob);
  return http({
    method: "post",
    url: "/api/v1/files",
    data: formData,
    headers: {
      "Content-Type": "multipart/form-data",
    },
  }).then(r => notifyPc(r, 'file'))
};
