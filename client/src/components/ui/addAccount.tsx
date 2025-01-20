import { Button } from './button'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from './dialog'

export const addAccount = (props : {}) => {
    return (
        <div>
 
<Dialog>

  <DialogContent>
    <DialogHeader>
      <DialogTitle>Are you absolutely sure?</DialogTitle>
      <DialogDescription>
        This action cannot be undone. This will permanently delete your account
        and remove your data from our servers. <br/>
            <Button variant="destructive" className="mt-5" >
            

  Im Sure!
  </Button>

      </DialogDescription>
    </DialogHeader>
  </DialogContent>
</Dialog>

        </div>
    )
}
