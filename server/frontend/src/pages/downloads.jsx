import { useQuery } from "../hooks/use_query";
import { BigTextarea, Button, Header, Layout } from "./home/components";
import styled from "styled-components";
import { Center } from "../components/center";
import { Space } from "../components/space";
import { useState } from "react";
import { useRef } from "react";
import { useEffect } from "react";
import { http } from "../shared/http";
import { history } from "../shared/history";

export const Downloads = () => {
  const query = useQuery()
  const type = normalizeType(query.type)
  const [text, setText] = useState("")
  useEffect(() => {
    if (type === "text") {
      http.get(query.url).then(({ data }) => {
        setText(data)
      })
    }
  }, [type])

  let node = null
  switch (type) {
    case 'text':
      node = (
        <div>
          <BigTextarea readOnly value={text} />
          <Space />
          <Center virtical>
            <Button>請手動複製上方文字</Button>
          </Center>
        </div>
      )
      break;
    case 'file':
      node = (
        <Center virtical>
          <a href={query.url}>
            <svg><use xlinkHref="#icon-file"></use></svg>
            <Space />
            <Center>
              <Button>點擊下載檔案</Button>
            </Center>
          </a>
        </Center>
      )
      break;
    case 'image':
      node = (
        <Center>
          <a href={query.url}>
            <Picture src={query.url} />
            <Center>
              <Button>長按或點擊，即可下載圖片</Button>
            </Center>
          </a>
        </Center>
      )
      break;
  }
  const onClickUpload = () => {
    history.push("/")
  }
  return (
    <Layout>
      <Header>同步傳</Header>
      {node}
      <Space x3 />
      <Center>
        <Button onClick={onClickUpload}>我也要上傳</Button>
      </Center>
    </Layout>
  )
};
const Picture = styled.img`
  border: 2px solid ${({ theme }) => theme.borderColor};
  margin: 16px;
`
const P = styled.p`
  margin: 8px 0;
`

const normalizeType = type => {
  if (/^image\/.*$/.test(type)) {
    return 'image'
  } else if (/^text(\/.*)?$/.test(type)) {
    return 'text'
  } else {
    return 'file'
  }
}