import Combine
import Foundation
import SwiftUI

class TaskListViewModel: ObservableObject {
    @Published var tasks: [ZTask]?
    @Published var openTasks: [ZTask]?
    @Published var todayOpenTasks: [ZTask]?
    @Published var isLoading: Bool = true

    @AppStorage("jwt") private var token: String?

    func loadTasks() {
        guard let token = token else {
            print("Token is missing")
            return
        }
        print("start")
        fetchTasks(token: token) { result in
            DispatchQueue.main.async {
                switch result {
                case .success(let fetchedTasks):
                    self.tasks = fetchedTasks
                case .failure(let error):
                    print(error)
                    print("Unable to load tasks: \(error.localizedDescription)")
                }
                self.isLoading = false
                print("done")
            }
        }
        print("loading tasks")
    }
}
