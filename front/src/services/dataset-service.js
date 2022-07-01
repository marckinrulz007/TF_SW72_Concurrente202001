import http from './http-common';

class DatasetService {

    getAllData(){
        return http.get(`data`)
    }
}

export default new DatasetService();
