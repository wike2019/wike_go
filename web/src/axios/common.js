import axios from 'axios';
import {ElMessage} from "element-plus";

// 创建 Axios 实例并配置
const apiClient = axios.create({
    baseURL: '/core/', // API 的基础 URL
    timeout: 30000, // 请求超时时间（毫秒）
});
apiClient.interceptors.response.use(
    function (response) {
        // 对响应数据做点什么
        if (response.data.code!=200){
            ElMessage.error(response.data.msg)
            return
        }
        return response.data;
    }, function (error) {
        // 对响应错误做点什么
        ElMessage.error(error+'系统异常')
    }
);


export {apiClient}